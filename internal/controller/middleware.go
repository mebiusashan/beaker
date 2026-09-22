package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
)

func LoginExpiredCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenKey, ok := sessionKey(c.GetHeader("X-Beaker-Session"), config.AuthEXPIRE_TIME); ok {
			c.Set("adminLoginKey", tokenKey)
			c.Next()
			return
		}
		ErrorFromCode(c, http.StatusUnauthorized, common.ErrorCodeUnauthorized, "Need Login")
		c.Abort()
		return
	}
}

func DecodeForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, err.Error())
			c.Abort()
			return
		}
		desKey := ""
		if value, ok := c.Get("adminLoginKey"); ok {
			desKey, _ = value.(string)
		}
		key, err := base64.StdEncoding.DecodeString(desKey)
		if err != nil {
			ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
			c.Abort()
			return
		}

		data64, err := cert.Base64Decode(string(data))
		if err != nil {
			ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
			c.Abort()
			return
		}

		sl, err := cert.TripleDesDecrypt(data64, key)
		if err != nil {
			ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(sl))
		var postData common.BaseReqMsg
		err = c.BindJSON(&postData)
		if err != nil {
			ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, err.Error())
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
			if err := controllerContext.Cache.ClearAll(); err != nil {
				ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeCache, err)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
