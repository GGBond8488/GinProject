package models

import (
	"My-gin-Project/pkg/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeleteOn int `json:"delete_on"`
}

func Setup()  {
	var (
		err error
		dbType, dbName, user, password, host, tablePrefix string
	)
	dbType = setting.DatabaseSetting.Type
	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tablePrefix = setting.DatabaseSetting.TablePrefix
	db,err = gorm.Open(dbType,fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&&parseTime=True&loc=Local",user,
		password,host,dbName))
	if err!=nil{
		log.Println(err)
	}
	//在所有模型的表的前面加上一个前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix+defaultTableName
	}
	// 全局禁用表名复数
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete",deleteCallback)
}

func CloseDB()  {
	defer db.Close()
}

//自定义回调
func updateTimeStampForCreateCallback(scop *gorm.Scope){
	if !scop.HasError(){
		nowTime := time.Now().Unix()
		if createTimeField,ok := scop.FieldByName("CreatedOn");ok{
			if createTimeField.IsBlank{
				createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField,ok := scop.FieldByName("ModifuiedOn");ok{
			if modifyTimeField.IsBlank{
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

//scope.Get("gorm:delete_option") 检查是否手动指定了delete_option
//scope.FieldByName("DeletedOn") 获取我们约定的删除字段，若存在则 UPDATE 软删除，若不存在则 DELETE 硬删除
//scope.QuotedTableName() 返回引用的表名，这个方法 GORM 会根据自身逻辑对表名进行一些处理
//scope.CombinedConditionSql() 返回组合好的条件SQL，看一下方法原型很明了

//scope.AddToVars 该方法可以添加值作为SQL的参数，也可用于防范SQL注入

func deleteCallback(scop *gorm.Scope)  {
	if !scop.HasError(){
		var extraOption string
		if str,ok := scop.Get("gorm:delete_option");ok{
			extraOption = fmt.Sprint(str)
		}

	deleteOnField,hasDeletedOnField := scop.FieldByName("DeleteOn")

	if !scop.Search.Unscoped && hasDeletedOnField{
		scop.Raw(fmt.Sprintf(
			"UPDATE %v SET %v=%v%v%v",
			scop.QuotedTableName(),
			scop.Quote(deleteOnField.DBName),
			scop.AddToVars(time.Now().Unix()),
			addExtraSpaceIfExist(scop.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOption),
			)).Exec()
	}else {
		scop.Raw(fmt.Sprintf(
			"DELETE FROM %v%v%v",
			scop.QuotedTableName(),
			addExtraSpaceIfExist(scop.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOption),
			)).Exec()
	}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}