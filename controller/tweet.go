package controller

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
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
	value, _ := c.Get("data")
	data := common.TweetModel{}
	json.Unmarshal(value.([]byte), &data)
	err := ct.Context.Model.TweetAdd(data.Content)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Tweet added successfully", nil)
}

func (ct *TweetController) Del(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.TweetModel{}
	json.Unmarshal(value.([]byte), &data)
	err := ct.Context.Model.TweetDel(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Tweet deleted successfully", nil)
}

func (ct *TweetController) List(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.TweetListResp{}
	json.Unmarshal(value.([]byte), &data)
	var page uint = 1
	if data.CurPage >= 1 {
		page = data.CurPage
	}

	twNums := ct.Context.Model.TweetCount()
	data.TweNum = twNums
	twNums = uint(math.Ceil(float64(twNums) / float64(10)))
	if page > twNums {
		page = twNums
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

	data.Code = 0
	data.CurPage = page
	data.TotlePage = twNums
	data.List = ts

	writeSucc(c, "", data)
}
