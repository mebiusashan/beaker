package cli

import (
	"fmt"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func ArtAll(host string, refresh bool, key []byte) {
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_ART_LIST, refresh, key, "")

	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-5d%s\n", uint(v["ID"].(float64)), v["Title"])
	}
}

func ArtRm(host string, refresh bool, key []byte, id uint) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	net.PostJsonWithEncrypt(host+net.CLI_ART_RM, refresh, key, sendData)
}

func ArtAdd(host string, refresh bool, key []byte, content string, title string, cid uint) {
	sendData := common.ArticleModel{Title: title, Content: content}
	sendData.Catid = cid
	net.PostJsonWithEncrypt(host+net.CLI_ART_ADD, refresh, key, sendData)
}

func ArtDownload(host string, refresh bool, key []byte, id uint) (string, string) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_ART_DOWN, refresh, key, sendData)
	data := jsonData.Data.(map[string]interface{})
	return data["Title"].(string), data["Content"].(string)
}

func ArtModify(host string, refresh bool, key []byte, id uint, catId uint, title string, content string) {

}
