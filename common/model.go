package common

import "github.com/jinzhu/gorm"

type BaseMsg struct {
	Code int `json:"code"`
	Data interface{}
}

type SuccMsg struct {
	BaseMsg
	Msg string `json:"msg"`
}

type LoginReq struct {
	DK string
	UN string
	PW string
}

type TweList struct {
	BaseMsg
	CurPage   uint
	TotlePage uint
	TweNum    uint
	List      []TweetDB
}
type TweetDB struct {
	gorm.Model
	Context string
}

type ArcDB struct {
	gorm.Model
	Catid   uint
	Title   string
	Context string
}

type CatDBDel struct {
	CatDB
	MvID uint
}
type CatDB struct {
	gorm.Model
	Cname string
	Name  string
}

type PageDB struct {
	gorm.Model
	Title   string
	Context string
}
