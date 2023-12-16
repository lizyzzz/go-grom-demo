package gormdemo

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
一对一关系比较少，一般用于表的扩展

例如一张用户表，有很多字段

那么就可以把它拆分为两张表，常用的字段放主表，不常用的字段放详情表
*/

type People struct {
	ID         uint
	Name       string
	Age        int
	Gender     bool
	PeopleInfo PeopleInfo // 通过UserInfo可以拿到用户详情信息
}

type PeopleInfo struct {
	PeopleID uint // 外键
	ID       uint
	Addr     string
	Like     string
}

// 创建关联表
func One2OneCreateTable() {
	DB.AutoMigrate(&People{}, &PeopleInfo{})
}

// 一对一添加
func One2OneInsertTable() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 新增一个用户
	DB.Create([]*People{
		{
			Name:   "资源",
			Age:    25,
			Gender: true,
			PeopleInfo: PeopleInfo{
				Addr: "湖南",
				Like: "code",
			},
		},
		{
			Name:       "自愿",
			Age:        24,
			Gender:     false,
			PeopleInfo: PeopleInfo{},
		},
	})

	// 新增附加信息
	DB.Create(&PeopleInfo{
		PeopleID: 2,
		Addr:     "南京市",
		Like:     "吃饭",
	})

}

func One2OneQuery() {
	var user People
	DB.Preload("PeopleInfo").Take(&user)
	fmt.Println(user)
}
