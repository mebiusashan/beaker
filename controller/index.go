package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type IndexController struct {
	BaseController
}

func (ct *IndexController) Do(c *gin.Context) {
	bodyStr, err := ct.Context.Cache.GET(common.TAG_HOME, "")
	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
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

	arcs, err := ct.Context.Model.ArticleFindWithNum(ct.Context.Config.Website.INDEX_LIST_NUM)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)
	vars.Set("ISINDEX", true)

	bodyStr, err = ct.Context.View.Render(common.TEMPLATE_HOME, vars)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	ct.Context.Cache.SETNX(common.TAG_HOME, "", bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}
