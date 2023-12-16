package gormdemo

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type User struct {
// 	ID       uint      `gorm:"size:4"`
// 	Name     string    `gorm:"size:8"`
// 	Articles []Article // 用户拥有的文章列表
// }

// type Article struct {
// 	ID     uint   `gorm:"size:4"`
// 	Title  string `gorm:"size:16"`
// 	UserID uint   `gorm:"size:4"` // 属于-默认使用<关联表名+ID>作为外键 (这里的类型要和引用的外键类型一致，包括大小)
// 	User   User   // 属于
// }

// 重写外键
type User struct {
	ID       uint      `gorm:"size:4"`
	Name     string    `gorm:"size:8;unique;not null"`
	Articles []Article `gorm:"foreignKey:UserID"` // 用户拥有的文章列表(这里的 <UserID> 同 Article 的 <UserID>)
}

type Article struct {
	ID     uint   `gorm:"size:4"`
	Title  string `gorm:"size:16;unique;not null"`
	UserID uint   `gorm:"size:4"`            // 属于(这里是外键字段)
	User   User   `gorm:"foreignKey:UserID"` // 这里的 <UserID> 的名字要与本类型的外键属性名相同
}

// 创建关联表
func One2MoreCreateTable() {
	DB.AutoMigrate(&User{}, &Article{})
}

// 一对多添加
func One2MoreInsertTable() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})
	// 1. 创建用户，并且创建文章
	a1 := Article{Title: "python"}
	a2 := Article{Title: "golang"}
	user := User{Name: "枫枫", Articles: []Article{a1, a2}}
	DB.Create(&user)
	a1 = Article{Title: "cpp"}
	a2 = Article{Title: "java"}
	user = User{Name: "资源", Articles: []Article{a1, a2}}
	DB.Create(&user)

	// 2. 创建文章，关联已有用户
	a3 := Article{Title: "golang零基础入门", UserID: 1}
	DB.Create(&a3)

	// 3. 给现有用户绑定文章
	var u User
	DB.Take(&u, 2)

	var article Article
	DB.Take(&article, 2)

	// u.Articles = []Article{article}
	// DB.Save(&u)

	// 用 append 方法(常用方法)
	DB.Model(&u).Association("Articles").Append(&article)

	// 4. 给现有文章关联用户
	u = User{}
	DB.Take(&u, 1)
	article = Article{}
	DB.Take(&article, 4)

	article.UserID = 1
	// DB.Save(&article)
	// 也可以用 append 方法(常用方法)
	DB.Model(&article).Association("User").Append(&u)
}

// 一对多查询
func One2MoreQuery() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 错误
	/*
		var user User
		DB.Take(&user, 1)
		fmt.Println(user)
	*/

	// 1. 预加载
	// 先预加载文章列表, 再查询用户
	var user User
	DB.Preload("Articles").Take(&user, 1)
	fmt.Println(user)
	// 先预加载用户列表, 再查询文章
	var article Article
	DB.Preload("User").Take(&article, 1)
	fmt.Println(article)

	// 2. 嵌套预加载
	// 查找文章所属用户的所有文章
	var a Article
	DB.Preload("User.Articles").Take(&a, 1)
	fmt.Println(a)

	// 3. 带条件的预加载
	var u User
	DB.Preload("Articles", "id = ?", 1).Take(&u, 1)
	fmt.Println(u)

	// 4. 自定义预加载
	user = User{}
	DB.Preload("Articles", func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", []int{1, 2})
	}).Take(&user, 1)
	fmt.Println(user)
}

// 一对多删除
func One2MoreDelete() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. (级联删除)删除用户, 与用户关联的文章也会删除
	var user User
	DB.Take(&user, 1)
	DB.Select("Articles").Delete(&user)

	// 2. (清除外键关系)将与用户关联的文章的外键设置为null
	user = User{}
	DB.Preload("Articles").Take(&user, 2)
	DB.Model(&user).Association("Articles").Delete(&user.Articles)
}
