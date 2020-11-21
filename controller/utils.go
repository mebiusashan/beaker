package controller

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/common"
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
