package controller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/russross/blackfriday"
)

func (ct *ArticleController) Do(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}
	if hasCacheWriteBody(c, ct.Context.Cache, common.TAG_ARCHIVE, idstr) {
		return
	}

	arts, err := ct.Context.Model.ArticleFindByID(uint(id))
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	markdownWithUnixLineEndings := strings.Replace(arts.Content, "\r\n", "\n", -1)
	bodyStr := string(blackfriday.MarkdownBasic([]byte(markdownWithUnixLineEndings)))
	vars := ct.Context.View.GetVarMap()
	vars.Set("body", bodyStr)
	vars.Set("title", arts.Title)

	bodyStr, err = ct.Context.View.Render(common.TEMPLATE_ARCHIVE, vars)
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}
	ct.Context.Cache.SETNX(common.TAG_ARCHIVE, idstr, bodyStr, ct.Context.Config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}
