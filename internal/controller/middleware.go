package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
)

func LoginExpiredCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetLoginInfo().CheckExpired(config.AuthEXPIRE_TIME) {
			writeFail(c, "Need Login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func DecodeForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			writeFail(c, err.Error())
			c.Abort()
			return
		}
		desKey := GetLoginInfo().loginKey
		key, err := base64.StdEncoding.DecodeString(desKey)
		if err != nil {
			writeFail(c, "No credit request")
			c.Abort()
			return
		}

		data64, err := cert.Base64Decode(string(data))
		if err != nil {
			writeFail(c, "No credit request")
			c.Abort()
			return
		}

		sl, err := cert.TripleDesDecrypt(data64, key)
		if err != nil {
			writeFail(c, "No credit request")
			c.Abort()
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(sl))
		var postData common.BaseReqMsg
		err = c.BindJSON(&postData)
		if err != nil {
			writeFail(c, err.Error())
			c.Abort()
			return
		}
		jsonString, _ := json.Marshal(postData.Data)
		c.Set("data", jsonString)
		c.Set("refresh", postData.Refresh)
		c.Next()
	}
}

func RefreshCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, has := c.Get("refresh")
		if has && value.(bool) {
			controllerContext.Cache.ClearAll()
		}
		c.Next()
	}
}
