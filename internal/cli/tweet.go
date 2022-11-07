package cli

import (
	"fmt"
	"strconv"

	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/net"
)

func TweetAll(host string, refresh bool, key []byte, curPage uint) {
	sendData := common.TweetListReq{CurPage: curPage}
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_TWEET_LIST, refresh, key, sendData)
	dd := jsonData.Data.(map[string]interface{})

	maxid := 0
	for _, va := range dd["List"].([]interface{}) {
		v := va.(map[string]interface{})
		l := len(strconv.Itoa(int(v["ID"].(float64))))
		if maxid < l {
			maxid = l
		}
	}
	for _, va := range dd["List"].([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-"+strconv.Itoa(maxid)+"d %s\n", uint(v["ID"].(float64)), v["Content"])
	}
	fmt.Println("---------------------------------------")
	fmt.Println(dd["TotlePage"], "pages,", dd["TweNum"], "tweets, current", dd["CurPage"], "page")
}

func TweetRm(host string, refresh bool, key []byte, id uint) {
	sendData := common.ArticleModel{}
	sendData.ID = id
	net.PostJsonWithEncrypt(host+net.CLI_TWEET_RM, refresh, key, sendData)
}

func TweetAdd(host string, refresh bool, key []byte, message string) {
	sendData := common.TweetModel{Content: message}
	net.PostJsonWithEncrypt(host+net.CLI_TWEET_ADD, refresh, key, sendData)
}
