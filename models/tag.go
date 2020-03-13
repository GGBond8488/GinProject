package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func GetTags(pageNum int,pageSize int,maps interface{})(tags []Tag,err error)  {
	err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return nil,err
	}
	return
}
func GetTagTotal(maps interface {}) (count int,err error){
	//model指定要运行数据库操作的模型
	// //将所有用户的名称更新为`hello'
	// db.Model（＆User {}）.Update（“ name”，“ hello”）
	// //如果用户的主键是非空白的，则将其用作条件，然后仅将用户名更新为`hello'
	// db.Model（＆user）.Update（“ name”，“ hello”）
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	return
}

func ExistTagByName(name string) (exists bool,err error) {
	var tag Tag
	err = db.Select("id").Where("name = ? AND delete_on = ?", name,0).First(&tag).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return false,err
	}
	if tag.ID > 0 {
		return true,nil
	}

	return false,nil
}

func AddTag(name string, state int, createdBy string)error{
	return db.Create(&Tag {
		Name : name,
		State : state,
		CreatedBy : createdBy,
	}).Error

}
func ExistTagByID(id int)(exists bool,err error) {
	var tag Tag
	db.Select("id").Where("id = ? AND delete_on = ?", id,0).First(&tag)
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return false,err
	}
	if tag.ID > 0 {
		return true,nil
	}

	return false,nil
}

func DeleteTag(id int) error {
	 return  db.Where("id = ?", id).Delete(&Tag{}).Error

}

func EditTag(id int, data interface {}) error {
	return db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error

}

func DeleteAllTag()error{
	return db.Unscoped().Where("delete_on != ?",0).Delete(&Tag{}).Error
}

//这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。
//
//gorm所支持的回调方法：
//
//创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
//更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
//删除：BeforeDelete、AfterDelete
//查询：AfterFind
//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}