package beaker

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (ct *ErrCtrl) Do404(c *gin.Context) {
	bodyStr, err := ct.ctrl.mvc.cache.GET(TAG_NOTFOUND, "")

	if err == nil && bodyStr != "" {
		c.Writer.WriteHeader(404)
		c.Writer.WriteString(bodyStr)
		return
	}

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("title", "404")
	str, err := ct.ctrl.mvc.view.Render(NotFound, vars)
	if err != nil {
		fmt.Println(err)
		ct.Do500(c)
		return
	}

	ct.ctrl.mvc.cache.SETNX(TAG_NOTFOUND, "", str, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
	c.Writer.WriteHeader(404)
	c.Writer.WriteString(str)
}

func (ct *ErrCtrl) Do500(c *gin.Context) {
	c.String(500, "500 Server Error")
}
