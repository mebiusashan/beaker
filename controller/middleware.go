package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/config"
)

func LoginExpiredCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetLoginInfo().CheckExpired(config.AuthEXPIRE_TIME) {
			writeFail(c, "Need Login")
			return
		}
		c.Next()
	}
}
