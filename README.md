# go-gorm-demo
go-grom 框架学习

### 1. gorm 连接 与 连接配置
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

### 2. 模型(相当于可以在创建表的时候指定字段属性)
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
### 3. 单表操作
#### 3.1 单表插入
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

#### 3.2 单表查询
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

#### 3.3 单表更新
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

#### 3.4 单表删除
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

### 4. 触发器
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

### 5. 高级查询
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


### 6. 一对一关系
```Go
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

// 一对一的查询
func One2OneQuery() {
	var user People
	DB.Preload("PeopleInfo").Take(&user)
	fmt.Println(user)
}

// 一对一的删除
func One2OneDelete() {
	var p People
	DB.Take(&p, 1)
	DB.Select("PeopleInfo").Delete(&p)
}
```

### 7. 一对多关系
#### 7.1 一对多的创建
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
#### 7.2 一对多的添加
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
#### 7.3 一对多的查询
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
#### 7.4 一对多的删除
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



### 8. 多对多关系
#### 8.1 多对多的创建
```Go
/*
需要用第三张表存储两张表的关系
*/

// 不同的文章有不同的标签

type Tag struct {
	ID     uint
	Name   string
	Papers []Paper `gorm:"many2many:article_tags;"` // 用于反向引用
}

type Paper struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags;"` // article_tags 是第三张表的名称
}

// article_tags 对应 ArticleTags, 以大写对应 `_`
type ArticleTags struct {
	PaperID   uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

// 多对多创建表
func More2MoreCreateTable() {
	DB.AutoMigrate(&Tag{}, &Paper{})
}
```
#### 8.2 多对多的添加
```Go
// 多对多添加
func More2MoreInsertTable() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. 添加文章, 并创建标签
	DB.Create(&Paper{
		Title: "python基础课程",
		Tags: []Tag{
			{Name: "python"},
			{Name: "基础课程"},
		},
	})

	// 2. 添加文章, 选择标签
	var tags []Tag
	DB.Find(&tags, "name = ?", "基础课程")
	DB.Create(&Paper{
		Title: "golang基础",
		Tags:  tags,
	})
}
```
#### 8.3 多对多的查询
```Go
// 多对多查询
func More2MoreQuery() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. 查询文章, 显示文章的标签列表
	var paper Paper
	DB.Preload("Tags").Take(&paper, 1)
	fmt.Println(paper)

	// 2. 查询标签, 显示文章列表
	var tag Tag
	DB.Preload("Papers").Take(&tag, 2)
	fmt.Println(tag)
}
```
#### 8.4 多对多的更新和删除
```Go
// 多对多更新
func More2MoreUpdate() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})
	// 1. 移除文章的标签
	var paper Paper
	DB.Preload("Tags").Take(&paper, 1)
	DB.Model(&paper).Association("Tags").Delete(paper.Tags)
	fmt.Println(paper)
	// 2. 更新文章的标签
	var p Paper
	var tags []Tag
	DB.Find(&tags, []int{2, 6, 7})

	DB.Preload("Tags").Take(&p, 2)
	DB.Model(&paper).Association("Tags").Replace(tags)
	fmt.Println(p)
}

// 多对多的删除
func More2MoreDelete() {
	var paper Paper
	DB.Take(&paper, 1)
	DB.Select("Tags").Delete(&paper)
}
```
#### 8.5 多对多自定义连接表
```Go
// 自定义连接表
func DefineConnTable() {
	DB.SetupJoinTable(&Paper{}, "Tags", &ArticleTags{}) // 设置新的连接表(因为两个类都有反向引用)
	DB.SetupJoinTable(&Tag{}, "Papers", &ArticleTags{}) // 设置新的连接表
	DB.AutoMigrate(&Paper{}, &Tag{}, &ArticleTags{})
}

// 自定义连接表的操作
func UserDefinedConnTableOP() {
	DB = DB.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 1. 添加文章并添加标签, 并自动关联
	DB.SetupJoinTable(&Paper{}, "Tags", &ArticleTags{}) // 要设置这个，才能走到我们自定义的连接表
	DB.SetupJoinTable(&Tag{}, "Papers", &ArticleTags{}) // 要设置这个，才能走到我们自定义的连接表
	DB.Create(&Paper{
		Title: "flask零基础入门",
		Tags: []Tag{
			{Name: "python"},
			{Name: "后端"},
			{Name: "web"},
		},
	})
	// CreatedAt time.Time 由于我们设置的是CreatedAt，gorm会自动填充当前时间，
	// 如果是其他的字段，需要使用到ArticleTag 的添加钩子 BeforeCreate

	// 2. 添加文章，关联已有标签
	var tags []Tag
	DB.Find(&tags, "name in ?", []string{"python", "web"})
	DB.Create(&Paper{
		Title: "flask请求对象",
		Tags:  tags,
	})

	// 3. 给已有文章关联标签
	paper := Paper{
		Title: "django基础",
	}
	DB.Create(&paper)
	var at Paper
	tags = []Tag{}
	DB.Find(&tags, "name in ?", []string{"python", "web"})
	DB.Take(&at, paper.ID).Association("Tags").Append(tags)

	// 4. 替换已有文章的标签
	var p Paper
	tags = []Tag{}
	DB.Find(&tags, "name in ?", []string{"后端"})
	DB.Take(&p, "title = ?", "django基础")
	DB.Model(&p).Association("Tags").Replace(tags)

	// 5. 查询文章列表，显示标签
	var papers []Paper
	DB.Preload("Tags").Find(&papers)
	fmt.Println(papers)
}
```

#### 8.6 多对多直接操作连接表
```Go
// 直接操作连接表
type UserModel struct {
	ID       uint
	Name     string
	Collects []ArticleModel `gorm:"many2many:user_collect_models;joinForeignKey:UserID;JoinReferences:ArticleID"`
}

type ArticleModel struct {
	ID    uint
	Title string
}

// UserCollectModel 用户收藏文章表
type UserCollectModels struct {
	UserID       uint         `gorm:"primaryKey"` // article_id
	UserModel    UserModel    `gorm:"foreignKey:UserID"`
	ArticleID    uint         `gorm:"primaryKey"` // tag_id
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID"`
	CreatedAt    time.Time
}

func QueryJoinTable() {
	DB.SetupJoinTable(&UserModel{}, "Collects", &UserCollectModels{}) // 设置新的连接表(因为两个类都有反向引用)
	DB.AutoMigrate(&UserModel{}, &ArticleModel{}, &UserCollectModels{})

	DB.Create(&UserModel{
		Name: "资源",
		Collects: []ArticleModel{
			{Title: "python"},
			{Title: "基础课程"},
		},
	})

	var collects []UserCollectModels

	var user UserModel
	DB.Take(&user, "name = ?", "资源")
	// 这里用map的原因是如果没查到，那就会查0值，如果是struct，则会忽略零值，全部查询
	DB.Preload("UserModel").Preload("ArticleModel").Where(map[string]any{"user_id": user.ID}).Find(&collects)

	for _, collect := range collects {
		fmt.Println(collect)
	}
}
```

### 9. 自定义数据类型
```Go
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
```

### 10. 事务
```Go
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
```
