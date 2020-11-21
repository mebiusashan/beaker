package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
	"github.com/russross/blackfriday"
)

type ArticleController struct {
	BaseController
}

func (ct *ArticleController) Do(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do404(c)
		return
	}

	bodyStr, err := ct.Context.Cache.GET(common.TAG_ARCHIVE, idstr)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	arts, err := ct.Context.Model.ArticleFindByID(uint(id))
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}
	bodyStr = string(blackfriday.Run([]byte(arts.Content)))

	vars := ct.Context.View.GetVarMap()
	vars.Set("body", bodyStr)
	vars.Set("title", arts.Title)

	bodyStr, err = ct.Context.View.Render(common.TEMPLATE_ARCHIVE, vars)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	ct.Context.Cache.SETNX(common.TAG_ARCHIVE, idstr, bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *ArticleController) Add(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.ArticleModel)
	err = ct.Context.Model.ArticleAdd(data.Catid, data.Title, data.Content)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Article added successfully", nil)
}

func (ct *ArticleController) Del(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.ArticleModel)
	err = ct.Context.Model.ArticleDel(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Article deleted successfully", nil)
}

func (ct *ArticleController) All(c *gin.Context) {
	arts, err := ct.Context.Model.ArticleFindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "All Article", arts)
}

func (ct *ArticleController) Down(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.ArticleModel)
	art, err := ct.Context.Model.ArticleFindByID(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}

	writeSucc(c, "Article", art)
}
