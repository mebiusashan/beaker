package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
	"github.com/russross/blackfriday"
)

type PageController struct {
	BaseController
}

func (ct *PageController) Do(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do404(c)
		return
	}

	bodyStr, err := ct.Context.Cache.GET(common.TAG_PAGE, idstr)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	page, err := ct.Context.Model.PageFindByID(uint(id))
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}
	bodyStr = string(blackfriday.Run([]byte(page.Content)))

	vars := ct.Context.View.GetVarMap()
	vars.Set("body", bodyStr)
	vars.Set("title", page.Title)

	bodyStr, err = ct.Context.View.Render(common.TEMPLATE_PAGE, vars)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	ct.Context.Cache.SETNX(common.TAG_PAGE, idstr, bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *PageController) Add(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.PageModel)
	err = ct.Context.Model.PageAdd(data.Title, data.Content)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Page added successfully", nil)
}

func (ct *PageController) Del(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.PageModel)
	err = ct.Context.Model.PageDel(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Page deleted successfully", nil)
}

func (ct *PageController) List(c *gin.Context) {
	pags, err := ct.Context.Model.PageFindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", pags)
}

func (ct *PageController) Down(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.PageModel)
	page, err := ct.Context.Model.PageFindByID(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", page)
}
