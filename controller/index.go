package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type IndexController struct {
	BaseController
}

func (ct *IndexController) Do(c *gin.Context) {
	if hasCacheWriteBody(c, ct.Context.Cache, common.TAG_HOME, "") {
		return
	}

	pages, err := ct.Context.Model.PageFindAll()
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	cats, err := ct.Context.Model.CategoryFindAll()
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}
	arcs, err := ct.Context.Model.ArticleFindWithNum(ct.Context.Config.Website.INDEX_LIST_NUM)
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)
	vars.Set("ISINDEX", true)

	bodyStr, err := ct.Context.View.Render(common.TEMPLATE_HOME, vars)
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	ct.Context.Cache.SETNX(common.TAG_HOME, "", bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}
