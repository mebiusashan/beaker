package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func TweetAll(host string, curPage uint) {
	postData := common.TweList{CurPage: curPage}
	jsonByte, err := json.Marshal(postData)
	common.Assert(err)

	jsonData := net.PostJson(host+"/admin/twe/list", strings.NewReader(string(jsonByte)))
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
		fmt.Printf("%-"+strconv.Itoa(maxid)+"d %s\n", uint(v["ID"].(float64)), v["Context"])
	}
	fmt.Println("---------------------------------------")
	fmt.Println(dd["TotlePage"], "pages,", dd["TweNum"], "tweets, current", dd["CurPage"], "page")
}
