package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type ArticleController struct {
	BaseController
}

func (ct *ArticleController) Add(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.ArticleModel{}
	json.Unmarshal(value.([]byte), &data)
	cat, err := ct.Context.Model.CategoryFindByID(data.Catid)
	if hasErrorWriteFail(c, err) {
		return
	}
	if cat.ID != data.Catid {
		writeFail(c, "category's id not found")
		return
	}
	md := writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err = ct.Context.Model.ArticleAdd(data.Catid, data.Title, md)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article added successfully", nil)
}

func (ct *ArticleController) Del(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.ArticleModel{}
	json.Unmarshal(value.([]byte), &data)
	err := ct.Context.Model.ArticleDel(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article deleted successfully", nil)
}

func (ct *ArticleController) All(c *gin.Context) {
	arts, err := ct.Context.Model.ArticleFindAll()
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "All Article", arts)
}

func (ct *ArticleController) Down(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.ArticleModel{}
	json.Unmarshal(value.([]byte), &data)
	art, err := ct.Context.Model.ArticleFindByID(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article", art)
}

func (ct *ArticleController) Modify(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.ArticleModel{}
	json.Unmarshal(value.([]byte), &data)
	data.Content = writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err := ct.Context.Model.ArticleUpdate(data.ID, &data)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article modify successfully", "")
}
