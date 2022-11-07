package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/config"
)

type OptionController struct {
	BaseController
}

func (ct *OptionController) Info(c *gin.Context) {
	path := ct.Context.Config.AuthInfo.ConfigPath
	config, err := config.NewWithPath(path, 0x1B)
	if hasErrorWriteFail(c, err) {
		return
	}
	writeSucc(c, "website", config.Website)
}

func (ct *OptionController) ClearCache(c *gin.Context) {
	ct.Context.Cache.ClearAll()
	writeSucc(c, "Clear cache successfully", nil)
}
