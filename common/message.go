package common

type BaseRespMsg struct {
	Code int `json:"code"`
	Data interface{}
}

type BaseReqMsg struct {
	refresh bool
}

type SuccMsgResp struct {
	BaseRespMsg
	Msg string `json:"msg"`
}

type LoginReq struct {
	BaseReqMsg
	DK string
	UN string
	PW string
}

type LoginResp struct {
	SuccMsgResp
}

type TweetListReq struct {
	BaseReqMsg
	CurPage uint
}

type TweetListResp struct {
	BaseRespMsg
	CurPage   uint
	TotlePage uint
	TweNum    uint
	List      []TweetModel
}
type CatRmReq struct {
	BaseReqMsg
	ID    uint
	Name  string
	Alias string
	MvID  uint
}
