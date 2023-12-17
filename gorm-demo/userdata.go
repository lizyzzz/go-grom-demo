package gormdemo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 自定义数据类型需要实现 Scan 和 Value 接口

// 1. 存放 json
type Info struct {
	Status string `json:"status"`
	Addr   string `json:"addr"`
	Age    int    `json:"age"`
}

// Scan 从数据库中读取出来(对象指针方法)
func (i *Info) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, i)
	return err
}

// Value 存入数据库(必须是对象方法)
func (i Info) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type AuthModel struct {
	ID   uint
	Name string
	Info Info `gorm:"type:string"`
}

func CreateUserDataTable() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	DB.AutoMigrate(&AuthModel{})

	// 插入数据
	DB.Create(&AuthModel{
		Name: "资源",
		Info: Info{
			Status: "alive",
			Addr:   "广东",
			Age:    25,
		},
	})

	// 查询数据
	var auth AuthModel
	DB.Take(&auth, "name = ?", "资源")
	fmt.Println(auth)

}
