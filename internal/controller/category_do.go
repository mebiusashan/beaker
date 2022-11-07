package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

func (ct *CategoryController) Do(c *gin.Context) {
	alias := c.Param("alias")
	if hasCacheWriteBody(c, ct.Context.Cache, common.TAG_CAT, alias) {
		return
	}

	pages, err := ct.Context.Model.PageFindAll()
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	cats, err := ct.Context.Model.CategoryFindAll()
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	cat, err := ct.Context.Model.CategoryFindByAlias(alias)
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	arcs, err := ct.Context.Model.ArticleFindByCatID(cat.ID)
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", cat.Name)
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)

	bodyStr, err := ct.Context.View.Render(common.TEMPLATE_CAT, vars)
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	ct.Context.Cache.SETNX(common.TAG_CAT, alias, bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}
