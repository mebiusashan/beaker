package controller

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/cache"
	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/config"
)

func write200(c *gin.Context, body string) {
	c.Writer.WriteHeader(200)
	c.Writer.WriteString(body)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeSucc(c *gin.Context, msg string, d interface{}) {
	data := new(common.SuccMsgResp)
	data.Code = common.SUCC
	data.Msg = msg
	data.Data = d
	c.JSON(200, data)
}

func writeFail(c *gin.Context, msg string) {
	data := new(common.SuccMsgResp)
	data.Code = common.FAIL
	data.Msg = msg
	c.JSON(200, data)
}

func hasErrorWriteFail(c *gin.Context, err error) bool {
	if err != nil {
		writeFail(c, err.Error())
		return true
	}
	return false
}

func writeStrSucc(c *gin.Context, msg string, d interface{}) {
	data := new(common.SuccMsgResp)
	data.Code = common.SUCC
	data.Msg = msg
	data.Data = d
	str, _ := json.Marshal(data)
	c.String(200, string((str)))
}

func writeStrFail(c *gin.Context, msg string) {
	data := new(common.SuccMsgResp)
	data.Code = common.FAIL
	data.Msg = msg
	str, _ := json.Marshal(data)
	c.String(200, string((str)))
}

func hasCacheWriteBody(c *gin.Context, cache *cache.Cache, tag string, key string) bool {
	bodyStr, err := cache.GET(tag, key)
	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return true
	}
	return false
}
func hasErrDo404(c *gin.Context, ct *ErrerController, err error) bool {
	if err != nil {
		ct.Do404(c)
		return true
	}
	return false
}

func hasErrDo500(c *gin.Context, ct *ErrerController, err error) bool {
	if err != nil {
		ct.Do500(c)
		return true
	}
	return false
}

func writeMarkdownImage(config config.Server, markdown string, imgs []common.ImgInfo) string {
	for i := 0; i < len(imgs); i++ {
		info := imgs[i]
		markdown = strings.ReplaceAll(markdown, info.Md5, config.SITE_URL+"/"+info.Md5)
		dec, err := base64.StdEncoding.DecodeString(info.Base64)
		if err != nil {
			continue
		}
		f, err := os.Create(config.STATIC_FILE_FOLDER + info.Md5 + "." + info.Suffix)
		if err != nil {
			continue
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			continue
		}
		if err := f.Sync(); err != nil {
			continue
		}
	}
	return markdown
}
