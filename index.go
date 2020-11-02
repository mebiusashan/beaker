package beaker

import (
	"github.com/gin-gonic/gin"
)

func (ct *IndCtrl) Do(c *gin.Context) {
	bodyStr, err := ct.ctrl.mvc.cache.GET(TAG_HOME, "")
	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	pages, err := ct.ctrl.mvc.model.PagDao.FindAll()
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	cats, err := ct.ctrl.mvc.model.CatDao.FindAll()
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	arcs, err := ct.ctrl.mvc.model.ArcDao.FindsByNum(ct.ctrl.mvc.config.Website.INDEX_LIST_NUM)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)
	vars.Set("ISINDEX", true)

	bodyStr, err = ct.ctrl.mvc.view.Render(HOME, vars)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	ct.ctrl.mvc.cache.SETNX(TAG_HOME, "", bodyStr, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}
