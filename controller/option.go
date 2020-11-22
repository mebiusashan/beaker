package controller

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/config"
)

type OptionController struct {
	BaseController
}

func (ct *OptionController) Info(c *gin.Context) {
	path := ct.Context.Config.AuthInfo.ConfigPath
	config, err := config.NewWithPath(path, 0x1B)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "website", config.Website)
}

func (ct *OptionController) ClearCache(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("新的", string(data))
	ct.Context.Cache.ClearAll()
	writeSucc(c, "Clear cache successfully", nil)
}
