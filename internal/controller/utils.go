package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/cache"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
	"github.com/russross/blackfriday"
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
	data.RequestID = requestID(c)
	data.Msg = msg
	data.Data = d
	c.JSON(http.StatusOK, data)
}

func writeFail(c *gin.Context, msg string) {
	writeFailCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, msg)
}

func writeFailCode(c *gin.Context, status int, errorCode string, msg string) {
	data := new(common.SuccMsgResp)
	data.Code = common.FAIL
	data.ErrorCode = errorCode
	data.RequestID = requestID(c)
	data.Msg = msg
	c.JSON(status, data)
}

func hasErrorWriteFail(c *gin.Context, err error) bool {
	if err != nil {
		writeFailCode(c, http.StatusInternalServerError, common.ErrorCodeDatabase, err.Error())
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
	data.ErrorCode = common.ErrorCodeInvalidRequest
	data.RequestID = requestID(nil)
	data.Msg = msg
	str, _ := json.Marshal(data)
	c.String(200, string((str)))
}

func requestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if value, ok := c.Get("requestID"); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}

func decodeAdminData(c *gin.Context, out interface{}) bool {
	value, ok := c.Get("data")
	if !ok {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, "missing request data")
		return false
	}
	data, ok := value.([]byte)
	if !ok {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, "invalid request data")
		return false
	}
	if err := json.Unmarshal(data, out); err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, err.Error())
		return false
	}
	return true
}

func renderMarkdown(content string) string {
	markdownWithUnixLineEndings := strings.ReplaceAll(content, "\r\n", "\n")
	renderer := blackfriday.HtmlRenderer(
		blackfriday.HTML_SKIP_HTML|
			blackfriday.HTML_SKIP_STYLE|
			blackfriday.HTML_SAFELINK|
			blackfriday.HTML_NOFOLLOW_LINKS|
			blackfriday.HTML_NOREFERRER_LINKS,
		"",
		"",
	)
	return string(blackfriday.Markdown([]byte(markdownWithUnixLineEndings), renderer, 0))
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
		ErrorFromCode(c, http.StatusNotFound, common.ErrorCodeNotFound, "resource not found")
		return true
	}
	return false
}

func hasErrDo500(c *gin.Context, ct *ErrerController, err error) bool {
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeInternal, err)
		return true
	}
	return false
}

func writeMarkdownImage(config config.Server, content string, imgs []common.ImgInfo) string {
	markdown := content
	for i := 0; i < len(imgs); i++ {
		info := imgs[i]
		if !validImageName(info.Md5, info.Suffix) {
			continue
		}
		markdown = strings.ReplaceAll(markdown, "("+info.Md5+info.Suffix, "("+config.SITE_URL+"/static/"+info.Md5+info.Suffix)
		dec, err := base64.StdEncoding.DecodeString(info.Base64)
		if err != nil {
			continue
		}
		if len(dec) == 0 || len(dec) > 5*1024*1024 {
			continue
		}
		if detected := http.DetectContentType(dec); !validImageContentType(detected, info.Suffix) {
			continue
		}
		f, err := os.CreateTemp(config.STATIC_FILE_FOLDER, ".beaker-upload-*")
		if err != nil {
			continue
		}
		if _, err := f.Write(dec); err != nil {
			f.Close()
			os.Remove(f.Name())
			continue
		}
		if err := f.Close(); err != nil {
			os.Remove(f.Name())
			continue
		}
		if err := os.Rename(f.Name(), filepath.Join(config.STATIC_FILE_FOLDER, info.Md5+info.Suffix)); err != nil {
			os.Remove(f.Name())
		}
	}
	return markdown
}

func validImageName(md5sum, suffix string) bool {
	if len(md5sum) != 32 || strings.Trim(md5sum, "0123456789abcdefABCDEF") != "" {
		return false
	}
	switch strings.ToLower(suffix) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return filepath.Base(md5sum+suffix) == md5sum+suffix
	default:
		return false
	}
}

func validImageContentType(contentType, suffix string) bool {
	switch strings.ToLower(suffix) {
	case ".jpg", ".jpeg":
		return contentType == "image/jpeg"
	case ".png":
		return contentType == "image/png"
	case ".gif":
		return contentType == "image/gif"
	case ".webp":
		return contentType == "image/webp"
	default:
		return false
	}
}
