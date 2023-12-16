package gormdemo

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// 日志显示
	// mysqlLogger := logger.Default.LogMode(logger.Info)

	// 连接 mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 跳过事务
		NamingStrategy: schema.NamingStrategy{ // 命名策略
			// TablePrefix:   "f_",  // 表格命名前缀
			SingularTable: true,  // 是否单数表名(不要小写)
			NoLowerCase:   false, // 不要自动转换大小写
		},
		// Logger: mysqlLogger, // 添加日志
	})

	if err != nil {
		panic("connect failed: " + err.Error())
	}
	DB = db
	fmt.Println(db)
}

type Student struct {
	ID     uint    `gorm:"size:10"`
	Name   string  `gorm:"size:16;not null"`
	Age    int32   `gorm:"size:3;not null"`
	Email  *string `gorm:"size:128"`
	Gender bool    `gorm:"default:false"`
}

type StudentInfo struct {
	Email  *string `gorm:"size:128"`
	Addr   string  `gorm:"type:varchar(32);not null"`
	Gender bool    `gorm:"default:false"`
}

// 触发器
// func (user *Student) BeforeCreate(tx *gorm.DB) error {
// 	email := "test@qq.com"
// 	user.Email = &email
// 	return nil
// }

func TestHook() {
	DB.Create(&Student{
		Name:   "lizy80",
		Age:    33,
		Gender: true,
	})
}

// 创建表操作
func CreateTable() {
	// 创建表
	DB.AutoMigrate(&Student{})     // 只新增表, 不删除表, 不修改表(字段大小会修改)
	DB.AutoMigrate(&StudentInfo{}) // 只新增, 不删除, 不修改
}

func ptrString(email string) *string {
	return &email
}

func PrepareData() {
	var stuList []Student
	DB.Find(&stuList).Delete(&stuList)
	stuList = []Student{
		{ID: 1, Name: "李元芳", Age: 32, Email: ptrString("lyf@yf.com"), Gender: true},
		{ID: 2, Name: "张武", Age: 18, Email: ptrString("zhangwu@lly.cn"), Gender: true},
		{ID: 3, Name: "枫枫", Age: 23, Email: ptrString("ff@yahoo.com"), Gender: true},
		{ID: 4, Name: "刘大", Age: 54, Email: ptrString("liuda@qq.com"), Gender: true},
		{ID: 5, Name: "李武", Age: 23, Email: ptrString("liwu@lly.cn"), Gender: true},
		{ID: 6, Name: "李琦", Age: 14, Email: ptrString("liqi@lly.cn"), Gender: false},
		{ID: 7, Name: "晓梅", Age: 25, Email: ptrString("xiaomeo@sl.com"), Gender: false},
		{ID: 8, Name: "如燕", Age: 26, Email: ptrString("ruyan@yf.com"), Gender: false},
		{ID: 9, Name: "魔灵", Age: 21, Email: ptrString("moling@sl.com"), Gender: true},
		{ID: 10, Name: "子源", Age: 25, Email: ptrString("ziyuan@zz.com"), Gender: true},
	}
	DB.Create(&stuList)
}

// 单表插入
func InsertTable() {
	// 添加记录
	email := "zzz@163.com"
	s1 := &Student{
		Name:   "lizy69",
		Age:    25,
		Email:  &email,
		Gender: false,
	}
	fmt.Println(s1.Gender)

	students := []*Student{
		{
			Name:   "lizy70",
			Age:    20,
			Gender: true,
		},
		{
			Name:   "lizy71",
			Age:    21,
			Email:  &email,
			Gender: false,
		},
	}

	// 插入单个
	tx := DB.Create(s1)
	// 批量插入
	tx = DB.Create(students)
	fmt.Println(tx)
}

// 单表查询
func QueryTable() {
	var student Student
	// 新建会话并添加日志
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. 返回查询的第一条记录, 不经过排序, 后续是查询条件
	tx := DB.Take(&student, "id = ?", 3).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}

	// 根据非主键查询
	student = Student{}
	tx = DB.Take(&student, fmt.Sprintf("email = '%s'", "zzz@163.com")).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}

	// 查询结果根据主键排序后的第一条
	student = Student{}
	tx = DB.First(&student, fmt.Sprintf("email = '%s'", "zzz@163.com")).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}
	// 查询结果根据主键排序后的最后一条
	student = Student{}
	tx = DB.Last(&student, fmt.Sprintf("email = '%s'", "zzz@163.com")).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}
	// 根据变量内置主键查询
	stu := Student{
		ID: 5,
	}
	tx = DB.Take(&stu).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(stu)
	}

	// 2. 查询多条记录
	// 根据条件查找
	var stuList []Student
	tx = DB.Find(&stuList, fmt.Sprintf("email = '%s'", "zzz@163.com")).Error
	if tx != gorm.ErrRecordNotFound {
		for _, s := range stuList {
			fmt.Println(s.ID, *s.Email)
		}
	}

	// 根据主键列表查询
	stuList = []Student{}
	tx = DB.Find(&stuList, []int{1, 3, 5, 6}).Error
	if tx != gorm.ErrRecordNotFound {
		for _, s := range stuList {
			fmt.Println(s.ID, s.Name)
		}
	}

	// 根据非主键列表查询
	stuList = []Student{}
	tx = DB.Select("id", "name").Find(&stuList, "name in ?", []string{"lizy66", "lizy70"}).Error
	if tx != gorm.ErrRecordNotFound {
		for _, s := range stuList {
			fmt.Println(s.ID, s.Name)
		}
	}
}

// 单表更新
func UpdateTable() {
	// 更新的前提是查找到数据

	var student Student
	email := "aaa@qq.com"
	email20 := "is20@qq.com"
	// 新建会话并添加日志
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. 单行更新单列
	// 更新全部字段
	tx := DB.Take(&student, 3).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}
	student.Age = 30
	DB.Save(&student)

	// 更新某个字段
	student = Student{}
	tx = DB.Take(&student, 5).Error
	if tx != gorm.ErrRecordNotFound {
		fmt.Println(student)
	}
	student.Age = 30
	student.Email = &email
	DB.Select("email").Save(&student)

	// 2. 批量更新
	var stuList []Student
	DB.Find(&stuList, "age = ?", 25).Update("email", "is25@qq.com")
	// 更简单的方式
	DB.Model(&Student{}).Where("age = ?", 21).Update("email", "is21@qq.com")

	// 更新多列
	DB.Model(Student{}).Where("age = ?", 20).Select("gender", "email").Updates(
		Student{
			Email:  &email20,
			Gender: true,
		},
	)
}

// 单表删除
func DeleteTable() {
	student := Student{
		ID: 4,
	}
	DB.Delete(&student)

	DB.Delete(&Student{}, []int{3, 5})

	// 基于切片删除
	// var stuList []Student
	// DB.Delete(&stuList)
}

// 高级查询
type User struct {
	Name2 string `gorm:"column:name"`
	Age   int
}

type AggeGroup struct {
	Gender int
	Count  int    `gorm:"column:count(id)"`
	Name   string `gorm:"column:group_concat(name)"`
}

func AdvancedQuery() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	var users []Student
	// 1.使用 where 查询
	DB.Where("name = ?", "子源").Find(&users)
	fmt.Println(users)

	DB.Where("name <> ?", "子源").Find(&users)
	fmt.Println(users)

	DB.Where("name in ?", []string{"子源", "张武"}).Find(&users)
	fmt.Println(users)

	DB.Where("name like ?", "李%").Find(&users)
	fmt.Println(users)

	DB.Where("age > ? and email like ?", "23", "%@qq.com").Find(&users)
	fmt.Println(users)

	DB.Where("gender = ? or email like ?", false, "%@qq.com").Find(&users)
	fmt.Println(users)

	// 2. 使用 结构体 查询
	DB.Where(&Student{Name: "李元芳", Age: 0}).Find(&users) // 会过滤默认值(零值)
	fmt.Println(users)

	// 3. 使用 map 查询
	DB.Where(map[string]any{"name": "李元芳", "age": 0}).Find(&users) // 不会过滤默认值(零值)
	fmt.Println(users)

	// 4. 使用 not 查询
	DB.Not("age > 23").Find(&users)
	fmt.Println(users)

	// 5. 使用 or 查询
	DB.Or("gender = ?", false).Or(" email like ?", "%@qq.com").Find(&users)
	fmt.Println(users)

	// 6. 使用 select 指定字段
	DB.Select("name", "age").Find(&users)
	fmt.Println(users)

	// 指定到另一个结构体
	var u []User
	DB.Select("name", "age").Find(&users).Scan(&u) // 查询两次(Find查找一次, Scan 查找一次), 去掉 Find 不行, 因为不知道表名
	fmt.Println(u)

	DB.Model(&Student{}).Select("name", "age").Scan(&u) // 查询一次(指定表名)
	fmt.Println(u)
	DB.Table("student").Select("name", "age").Scan(&users) // 查询一次(指定表名)
	fmt.Println(users)

	// 7. 排序
	DB.Order("age desc").Find(&users)
	fmt.Println(users)

	// 8. 分页
	// 一页两条，第1页
	DB.Limit(2).Offset(0).Find(&users)
	fmt.Println(users)
	// 第2页
	DB.Limit(2).Offset(2).Find(&users)
	fmt.Println(users)
	// 第3页
	DB.Limit(2).Offset(4).Find(&users)
	fmt.Println(users)

	// 9. 去重
	var ageList []int
	DB.Model(&Student{}).Select("age").Distinct("age").Scan(&ageList)
	fmt.Println(ageList)

	// 10. 分组查询
	DB.Model(&Student{}).Select("count(id)").Group("gender").Scan(&ageList)
	fmt.Println(ageList)

	var agge []AggeGroup
	DB.Model(&Student{}).Select("count(id)", "gender", "group_concat(name)").Group("gender").Scan(&agge)
	fmt.Println(agge)

	// 11. 原生 SQL
	DB.Raw(`SELECT count(id), gender, group_concat(name) FROM student GROUP BY gender`).Scan(&agge)
	fmt.Println(agge)

	// 12. 子查询
	// select * from students where age > (select avg(age) from students)
	DB.Model(Student{}).Where("age > (?)", DB.Model(Student{}).Select("avg(age)")).Find(&users) // 里面不用再写一个 Find
	fmt.Println(users)

	// 13. 命令参数(具名参数)
	DB.Where("name = @name and age = @age", sql.Named("name", "枫枫"), sql.Named("age", 23)).Find(&users)
	DB.Where("name = @name and age = @age", map[string]any{"name": "枫枫", "age": 23}).Find(&users)
	fmt.Println(users)

	// 14. find 的数据报存在 map
	var res []map[string]any
	DB.Model(Student{}).Find(&res)
	fmt.Println(res)

	// 15. 查询引用 scope
	DB.Scopes(Age23).Find(&users)
	fmt.Println(users)

}

// 可以在 model 层写一些通用的查询方式，这样外界就可以直接调用方法即可
func Age23(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 23)
}
