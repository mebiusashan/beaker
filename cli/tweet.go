package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/common"
)

func TweetAll(host string, curPage uint) {
	postData := common.TweList{CurPage: curPage}
	jsonByte, err := json.Marshal(postData)
	common.Assert(err)

	resp, err := http.Post(host+"/admin/twe/list", "", strings.NewReader(string(jsonByte)))
	common.Assert(err)
	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
		return
	}

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
