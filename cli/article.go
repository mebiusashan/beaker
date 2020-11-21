package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func ArtAll(host string) {
	jsonData := net.PostJson(host+net.CLI_ART_LIST, strings.NewReader(""))

	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-5d%s\n", uint(v["ID"].(float64)), v["Title"])
	}
}

func ArtRm(host string, id uint) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_ART_RM, strings.NewReader(string(jsonByte)))
}

func ArtAdd(host string, content string, title string, cid uint) {
	sendData := common.ArticleModel{Title: title, Content: content}
	sendData.Catid = cid
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_ART_ADD, strings.NewReader(string(jsonByte)))
}

func ArtDownload(host string, id uint) (string, string) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	jsonData := net.PostJson(host+net.CLI_ART_DOWN, strings.NewReader(string(jsonByte)))
	data := jsonData.Data.(map[string]interface{})
	return data["Title"].(string), data["Content"].(string)
}

func ArtModify(host string, id uint, catId uint, title string, content string) {

}
