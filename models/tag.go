package models

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func GetTags(pageNum int,pageSize int,maps interface{})(tags []Tag,ok bool)  {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	ok = true
	return
}
func GetTagTotal(maps interface {}) (count int,ok bool){
	//model指定您要运行数据库操作的模型
	// //将所有用户的名称更新为`hello'
	// db.Model（＆User {}）.Update（“ name”，“ hello”）
	// //如果用户的主键是非空白的，则将其用作条件，然后仅将用户名更新为`hello'
	// db.Model（＆user）.Update（“ name”，“ hello”）
	db.Model(&Tag{}).Where(maps).Count(&count)
	ok = true
	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ? AND delete_on = ?", name,0).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func AddTag(name string, state int, createdBy string) bool{
	db.Create(&Tag {
		Name : name,
		State : state,
		CreatedBy : createdBy,
	})

	return true
}
func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ? AND delete_on = ?", id,0).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface {}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
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