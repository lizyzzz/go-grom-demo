package gormdemo

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	username = "lizy"
	password = "123456"
	host     = "127.0.0.1"
	port     = 3306
	dbName   = "webserver"
	timeout  = "10s"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, dbName, timeout)

	// 连接 mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 跳过事务
		NamingStrategy: schema.NamingStrategy{ // 命名策略
			TablePrefix:   "f_",  // 表格命名前缀
			SingularTable: true,  // 是否单数表名(不要小写)
			NoLowerCase:   false, // 不要自动转换大小写
		},
	})

	if err != nil {
		panic("connect failed: " + err.Error())
	}
	DB = db
	fmt.Println(db)
}

type Student struct {
	ID   uint
	Name string
	Age  int32
}

func CreateTable() {
	// 创建表
	DB.AutoMigrate(&Student{})
}

func Connect() {

}
