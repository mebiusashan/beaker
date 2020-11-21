package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
)

type ErrerController struct {
	BaseController
}

func (ct *ErrerController) Do404(c *gin.Context) {
	bodyStr, err := ct.Context.Cache.GET(common.TEMPLATE_NotFound, "")

	if err == nil && bodyStr != "" {
		c.Writer.WriteHeader(404)
		c.Writer.WriteString(bodyStr)
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", "404")
	str, err := ct.Context.View.Render(common.TEMPLATE_NotFound, vars)
	if err != nil {
		fmt.Println(err)
		ct.Do500(c)
		return
	}

	ct.Context.Cache.SETNX(common.TAG_NOTFOUND, "", str, ct.Context.Config.Redis.EXPIRE_TIME)
	c.Writer.WriteHeader(404)
	c.Writer.WriteString(str)
}

func (ct *ErrerController) Do500(c *gin.Context) {
	c.String(500, "500 Server Error")
}
