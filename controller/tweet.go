package controller

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type TweetController struct {
	BaseController
}

func (ct *TweetController) Do(c *gin.Context) {
	page := ct.convParam(c.Param("page"))
	numOfOnePage := ct.Context.Config.Website.TWEET_NUM_ONE_PAGE

	twCount := ct.Context.Model.TweetCount()
	twCount = uint(math.Ceil(float64(twCount) / float64(numOfOnePage)))
	if page > twCount {
		page = twCount
	}

	twpages := ct.createPageNum(twCount)

	tws, err := ct.Context.Model.TweetFindByNum(page, numOfOnePage)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	pages, err := ct.Context.Model.PageFindAll()
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	cats, err := ct.Context.Model.CategoryFindAll()
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("tws", tws)
	vars.Set("twpages", twpages)
	vars.Set("twcurpage", page)

	bodyStr, err := ct.Context.View.Render(common.TEMPLATE_TWEET, vars)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}
	write200(c, bodyStr)
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
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.TweetModel)
	err = ct.Context.Model.TweetAdd(data.Content)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Tweet added successfully", nil)
}

func (ct *TweetController) Del(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	data := postData.Data.(common.TweetModel)
	err = ct.Context.Model.TweetDel(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Tweet deleted successfully", nil)
}

func (ct *TweetController) List(c *gin.Context) {
	var data common.BaseReqMsg
	err := c.BindJSON(&data)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	postData := data.Data.(common.TweetListResp)
	var page uint = 1
	if postData.CurPage >= 1 {
		page = postData.CurPage
	}

	twNums := ct.Context.Model.TweetCount()
	postData.TweNum = twNums
	twNums = uint(math.Ceil(float64(twNums) / float64(10)))
	if page > twNums {
		page = twNums
	}

	tws, err := ct.Context.Model.TweetFindByNum(page, 10)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	var ts []common.TweetModel
	for _, v := range tws {
		t := common.TweetModel{Content: v.Content}
		t.ID = v.ID
		ts = append(ts, t)
	}

	postData.Code = 0
	postData.CurPage = page
	postData.TotlePage = twNums
	postData.List = ts

	writeSucc(c, "", postData)
}
