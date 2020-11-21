package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type CategoryController struct {
	BaseController
}

func (ct *CategoryController) Do(c *gin.Context) {
	alias := c.Param("name")
	bodyStr, err := ct.Context.Cache.GET(common.TAG_CAT, alias)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	pages, err := ct.Context.Model.PageFindAll()
	if err != nil {
		fmt.Println(1, err)
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	cats, err := ct.Context.Model.CategoryFindAll()
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	cat, err := ct.Context.Model.CategoryFindByAlias(alias)
	if err != nil {
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	arcs, err := ct.Context.Model.ArticleFindByCatID(cat.ID)
	if err != nil {
		fmt.Println(3, err)
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", cat.Name)
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)

	bodyStr, err = ct.Context.View.Render(common.TEMPLATE_CAT, vars)
	if err != nil {
		fmt.Println(4, err)
		ct.Context.Ctrl.ErrC.Do500(c)
		return
	}

	ct.Context.Cache.SETNX(common.TAG_CAT, alias, bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *CategoryController) Add(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.CatModel)
	if data.Name == "" || data.Alias == "" {
		writeFail(c, "Null values ​​are not allowed")
		return
	}

	err = ct.Context.Model.CategoryAdd(data.Name, data.Alias)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Category added successfully", nil)
}

func (ct *CategoryController) Del(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.CatRmReq)
	err = ct.Context.Model.ArticleUpdateCat(data.ID, data.MvID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.Context.Model.CategoryDel(data.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Category deleted successfully", nil)
}

func (ct *CategoryController) All(c *gin.Context) {
	list, err := ct.Context.Model.CategoryFindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "Category list", list)
}

func (ct *CategoryController) Update(c *gin.Context) {
	var postData common.BaseReqMsg
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data := postData.Data.(common.CatModel)
	err = ct.Context.Model.CategoryUpdate(data.ID, data.Name, data.Alias)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if postData.Refresh {
		ct.Context.Cache.ClearAll()
	}
	writeSucc(c, "Category modified successfully", nil)
}
