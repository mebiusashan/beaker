package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func PageAll(host string) {
	jsonData := net.PostJson(host+net.CLI_PAGE_LIST, strings.NewReader(""))

	maxid := 0
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		l := len(strconv.Itoa(int(v["ID"].(float64))))
		if maxid < l {
			maxid = l
		}
	}
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-"+strconv.Itoa(maxid)+"d %s\n", uint(v["ID"].(float64)), v["Title"])
	}
}

func PageRm(host string, id uint) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_PAGE_RM, strings.NewReader(string(jsonByte)))
}

func PageAdd(host string, content string, title string) {
	sendData := common.PageModel{Title: title, Content: content}
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_PAGE_ADD, strings.NewReader(string(jsonByte)))
}

func PageDownload(host string, id uint) (string, string) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	jsonData := net.PostJson(host+net.CLI_PAGE_DOWN, strings.NewReader(string(jsonByte)))
	data := jsonData.Data.(map[string]interface{})
	return data["Title"].(string), data["Content"].(string)
}

func PageModify(host string, id uint, title string, content string) {

}
