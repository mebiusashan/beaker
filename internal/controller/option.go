package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
)

type OptionController struct {
	BaseController
}

func (ct *OptionController) Info(c *gin.Context) {
	path := ct.Context.Config.AuthInfo.ConfigPath
	config, err := config.NewWithPath(path, 0x1B)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeInternal, err)
		return
	}
	writeSucc(c, "website", config.Website)
}

func (ct *OptionController) ClearCache(c *gin.Context) {
	if err := ct.Context.Cache.ClearAll(); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeCache, err)
		return
	}
	writeSucc(c, "Clear cache successfully", nil)
}
