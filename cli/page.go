package cli

import (
	"fmt"
	"strconv"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func PageAll(host string, refresh bool, key []byte) {
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_PAGE_LIST, refresh, key, "")

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

func PageRm(host string, refresh bool, key []byte, id uint) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	net.PostJsonWithEncrypt(host+net.CLI_PAGE_RM, refresh, key, sendData)
}

func PageAdd(host string, refresh bool, key []byte, content string, title string, imgInfos []common.ImgInfo) {
	sendData := common.PageModel{Title: title, Content: content, Imgs: imgInfos}
	net.PostJsonWithEncrypt(host+net.CLI_PAGE_ADD, refresh, key, sendData)
}

func PageDownload(host string, refresh bool, key []byte, id uint) (string, string) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_PAGE_DOWN, refresh, key, sendData)
	data := jsonData.Data.(map[string]interface{})
	return data["Title"].(string), data["Content"].(string)
}

func PageModify(host string, refresh bool, key []byte, id uint, title string, content string, imgInfos []common.ImgInfo) {
	sendData := common.PageModel{Title: title, Content: content, Imgs: imgInfos}
	sendData.ID = id
	net.PostJsonWithEncrypt(host+net.CLI_PAGE_MODIFY, refresh, key, sendData)
}
