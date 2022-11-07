package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

type CategoryController struct {
	BaseController
}

func (ct *CategoryController) Add(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.CatModel{}
	json.Unmarshal(value.([]byte), &data)
	if data.Name == "" || data.Alias == "" {
		writeFail(c, "Null values ​​are not allowed")
		return
	}
	if data.Alias == "static" {
		writeFail(c, "static is a reserved word")
		return
	}
	err := ct.Context.Model.CategoryAdd(data.Name, data.Alias)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Category added successfully", nil)
}

func (ct *CategoryController) Del(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.CatRmReq{}
	json.Unmarshal(value.([]byte), &data)
	mcat, err := ct.Context.Model.CategoryFindByID(data.MvID)
	if hasErrorWriteFail(c, err) {
		return
	}
	if mcat.ID != data.MvID {
		writeFail(c, "Target category ID does not exist")
		return
	}
	err = ct.Context.Model.ArticleUpdateCat(data.ID, data.MvID)
	if hasErrorWriteFail(c, err) {
		return
	}
	err = ct.Context.Model.CategoryDel(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Category deleted successfully", nil)
}

func (ct *CategoryController) All(c *gin.Context) {
	list, err := ct.Context.Model.CategoryFindAll()
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Category list", list)
}

func (ct *CategoryController) Update(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.CatModel{}
	json.Unmarshal(value.([]byte), &data)
	if data.Alias == "static" {
		writeFail(c, "static is a reserved word")
		return
	}
	err := ct.Context.Model.CategoryUpdate(data.ID, data.Name, data.Alias)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Category modified successfully", nil)
}
