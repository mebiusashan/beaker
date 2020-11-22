package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type ArticleController struct {
	BaseController
}

func (ct *ArticleController) Add(c *gin.Context) {
	value, _ := c.Get("data")
	data := value.(common.ArticleModel)
	err := ct.Context.Model.ArticleAdd(data.Catid, data.Title, data.Content)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article added successfully", nil)
}

func (ct *ArticleController) Del(c *gin.Context) {
	value, _ := c.Get("data")
	data := value.(common.ArticleModel)
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
	data := value.(common.ArticleModel)
	art, err := ct.Context.Model.ArticleFindByID(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Article", art)
}
