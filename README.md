# go-gorm-demo
go-grom 框架学习

### gorm 连接 与 连接配置
需要下载mysql的驱动
```bash
go get gorm.io/driver/mysql
go get gorm.io/gorm
```
```Go
// 日志显示
mysqlLogger := logger.Default.LogMode(logger.Info)

// 连接 mysql
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  SkipDefaultTransaction: true, // 跳过事务
  NamingStrategy: schema.NamingStrategy{ // 命名策略
    TablePrefix:   "f_",  // 表格命名前缀
    SingularTable: true,  // 是否单数表名(不要小写)
    NoLowerCase:   false, // 不要自动转换大小写
  },
  Logger: mysqlLogger, // 添加日志
})
```

### 模型(相当于可以在创建表的时候指定字段属性)
* `type` 定义字段类型  
* `size` 定义字段大小  
* `column` 自定义列名  
* `primaryKey` 将列定义为主键  
* `unique` 将列定义为唯一键  
* `default` 定义列的默认值  
* `not null` 不可为空  
* `embedded` 嵌套字段  
* `embeddedPrefix` 嵌套字段前缀  
* `comment` 注释  
* **多个标签之前用 ; 连接**  

### 单表插入
```Go
func InsertTable() {
	// 添加记录
	email := "zzz@163.com"
	M := true
	F := false
	s1 := &Student{
		Name:   "lizy66",
		Age:    25,
		Email:  &email,
		Gender: &F,
	}
	fmt.Println(s1.Gender)

	students := []*Student{
		{
			Name:   "lizy67",
			Age:    20,
			Gender: &M,
		},
		{
			Name:   "lizy68",
			Age:    21,
			Email:  &email,
			Gender: &F,
		},
	}

	// 插入单个
	tx := DB.Create(s1)
	// 批量插入
	tx = DB.Create(students)
	fmt.Println(tx)
}
```

### 单表查询
```Go
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
	tx = DB.Find(&stuList, "name in ?", []string{"lizy66", "lizy70"}).Error
	if tx != gorm.ErrRecordNotFound {
		for _, s := range stuList {
			fmt.Println(s.ID, s.Name)
		}
	}
}
```

### 单表更新
```Go
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
```

### 单表删除
```Go
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
```

### 触发器
```Go
// 触发器
func (user *Student) BeforeCreate(tx *gorm.DB) error {
	email := "test@qq.com"
	user.Email = &email
	return nil
}

func TestHook() {
	DB.Create(&Student{
		Name:   "lizy80",
		Age:    33,
		Gender: true,
	})
}
```

### 高级查询
```Go
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

	// 13. 命令参数
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
```

### 一对多
#### 一对多的创建
```Go
type User struct {
	ID       uint      `gorm:"size:4"`
	Name     string    `gorm:"size:8"`
	Articles []Article `gorm:"foreignKey:UID"` // 用户拥有的文章列表(这里的 <UID> 同 Article 的 <UID>)
}

type Article struct {
	ID    uint   `gorm:"size:4"`
	Title string `gorm:"size:16"`
	UID   uint   `gorm:"size:4"`         // 属于(这里是外键字段)
	User  User   `gorm:"foreignKey:UID"` // 这里的 <UID> 的名字要与本类型的外键属性名相同
}

func One2MoreCreateTable() {
	DB.AutoMigrate(&User{}, &Article{})
}
```
#### 一对多的添加
```Go
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
```
#### 一对多的查询
```Go
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
```
#### 一对多的删除
```Go
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
```
