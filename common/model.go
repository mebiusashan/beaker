package common

import "github.com/jinzhu/gorm"

type ArticleModel struct {
	gorm.Model
	Catid   uint
	Title   string
	Content string
}

type PageModel struct {
	gorm.Model
	Title   string
	Content string
}

type TweetModel struct {
	gorm.Model
	Content string
}

type CatModel struct {
	gorm.Model
	Name  string
	Alias string
}
