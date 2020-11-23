package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type PageController struct {
	BaseController
}

func (ct *PageController) Add(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.PageModel{}
	json.Unmarshal(value.([]byte), &data)
	md := writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err := ct.Context.Model.PageAdd(data.Title, md)
	if hasErrorWriteFail(c, err) {
		return
	}

	writeSucc(c, "Page added successfully", nil)
}
func (ct *PageController) Del(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.PageModel{}
	json.Unmarshal(value.([]byte), &data)
	err := ct.Context.Model.PageDel(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Page deleted successfully", nil)
}

func (ct *PageController) List(c *gin.Context) {
	pags, err := ct.Context.Model.PageFindAll()
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "", pags)
}

func (ct *PageController) Down(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.PageModel{}
	json.Unmarshal(value.([]byte), &data)
	page, err := ct.Context.Model.PageFindByID(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "", page)
}

func (ct *PageController) Modify(c *gin.Context) {
	value, _ := c.Get("data")
	data := common.PageModel{}
	json.Unmarshal(value.([]byte), &data)
	data.Content = writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err := ct.Context.Model.PageUpdate(data.ID, &data)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Page modify successfully", "")
}
