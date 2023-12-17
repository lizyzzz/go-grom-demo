package gormdemo

import (
	"fmt"

	"gorm.io/gorm"
)

type Consumer struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Money int    `json:"money"`
}

func CreateTxTable() {
	DB.Set("gorm:tale_option", "ENGINE=InnoDB").AutoMigrate(&Consumer{})
	DB.Create([]*Consumer{
		{
			Name:  "张三",
			Money: 1000,
		},
		{
			Name:  "李四",
			Money: 500,
		},
	})
}

func TxExample() {
	var zhangsan, lisi Consumer
	DB.Take(&zhangsan, "name = ?", "张三")
	DB.Take(&lisi, "name = ?", "李四")
	// 张三给李四转账100元
	DB.Transaction(func(tx *gorm.DB) error {

		// 先给张三-100
		zhangsan.Money -= 100
		err := tx.Model(&zhangsan).Update("money", zhangsan.Money).Error
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 再给李四+100
		lisi.Money += 100
		err = tx.Model(&lisi).Update("money", lisi.Money).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 提交事务
		return nil
	})

	// 手动的方式
	// // 张三给李四转账100元
	// tx := DB.Begin()

	// // 先给张三-100
	// zhangsan.Money -= 100
	// err := tx.Model(&zhangsan).Update("money", zhangsan.Money).Error
	// if err != nil {
	// 	tx.Rollback()
	// }

	// // 再给李四+100
	// lisi.Money += 100
	// err = tx.Model(&lisi).Update("money", lisi.Money).Error
	// if err != nil {
	// 	tx.Rollback()
	// }
	// // 提交事务
	// tx.Commit()
}
