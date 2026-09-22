package common

type BaseRespMsg struct {
	Code      int         `json:"code"`
	ErrorCode string      `json:"errorCode,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type BaseReqMsg struct {
	Refresh bool        `json:"refresh"`
	Data    interface{} `json:"data"`
}

type SuccMsgResp struct {
	BaseRespMsg
	Msg          string `json:"msg"`
	SessionToken string `json:"sessionToken,omitempty"`
}

type LoginReq struct {
	BaseReqMsg
	DK string `json:"dk"`
	UN string `json:"un"`
	PW string `json:"pw"`
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
