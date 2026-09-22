package controller

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

type TweetController struct {
	BaseController
}

func (ct *TweetController) convParam(param string) uint {
	if param == "/" {
		return 1
	}
	str := param[1:]
	num, err := strconv.Atoi(str)
	if err != nil {
		return 1
	}
	return uint(num)
}

func (ct *TweetController) createPageNum(count uint) []uint {
	twpages := make([]uint, count)
	var i uint = 1
	for ; i <= count; i++ {
		twpages[i-1] = i
	}
	return twpages
}

func (ct *TweetController) Add(c *gin.Context) {
	data := common.TweetModel{}
	if !decodeAdminData(c, &data) {
		return
	}
	if data.Content == "" {
		writeFail(c, "tweet content is required")
		return
	}
	err := ct.Context.Model.TweetAdd(data.Content)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Tweet added successfully", nil)
}

func (ct *TweetController) Del(c *gin.Context) {
	data := common.TweetModel{}
	if !decodeAdminData(c, &data) {
		return
	}
	err := ct.Context.Model.TweetDel(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Tweet deleted successfully", nil)
}

func (ct *TweetController) List(c *gin.Context) {
	data := common.TweetListResp{}
	if !decodeAdminData(c, &data) {
		return
	}
	var page uint = 1
	if data.CurPage >= 1 {
		page = data.CurPage
	}

	twNums, err := ct.Context.Model.TweetCount()
	if hasErrorWriteFail(c, err) {
		return
	}
	data.TweNum = twNums
	totalPages := uint(math.Ceil(float64(twNums) / float64(10)))
	if totalPages == 0 {
		data.Code = common.SUCC
		data.CurPage = 1
		data.TotlePage = 0
		data.List = []common.TweetModel{}
		writeSucc(c, "", data)
		return
	}
	if page > totalPages {
		page = totalPages
	}

	tws, err := ct.Context.Model.TweetFindByNum(page, 10)
	if hasErrorWriteFail(c, err) {
		return
	}

	var ts []common.TweetModel
	for _, v := range tws {
		t := common.TweetModel{Content: v.Content}
		t.ID = v.ID
		ts = append(ts, t)
	}

	data.Code = common.SUCC
	data.CurPage = page
	data.TotlePage = totalPages
	data.List = ts

	writeSucc(c, "", data)
}
