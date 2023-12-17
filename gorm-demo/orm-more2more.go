package gormdemo

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
