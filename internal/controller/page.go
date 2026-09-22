package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

type PageController struct {
	BaseController
}

func (ct *PageController) Add(c *gin.Context) {
	data := common.PageModel{}
	if !decodeAdminData(c, &data) {
		return
	}
	if data.Title == "" {
		writeFail(c, "title is required")
		return
	}
	md := writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err := ct.Context.Model.PageAdd(data.Title, md)
	if hasErrorWriteFail(c, err) {
		return
	}

	writeSucc(c, "Page added successfully", nil)
}
func (ct *PageController) Del(c *gin.Context) {
	data := common.PageModel{}
	if !decodeAdminData(c, &data) {
		return
	}
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
	data := common.PageModel{}
	if !decodeAdminData(c, &data) {
		return
	}
	page, err := ct.Context.Model.PageFindByID(data.ID)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "", page)
}

func (ct *PageController) Modify(c *gin.Context) {
	data := common.PageModel{}
	if !decodeAdminData(c, &data) {
		return
	}
	data.Content = writeMarkdownImage(ct.Context.Config.Server, data.Content, data.Imgs)
	err := ct.Context.Model.PageUpdate(data.ID, &data)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "Page modify successfully", "")
}
