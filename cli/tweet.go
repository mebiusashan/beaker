package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/xlab/tablewriter"
)

func TweetAll(host string) {
	postData := common.TweList{CurPage: 0}
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

	tablewriter.EnableUTF8()
	table := tablewriter.CreateTable()
	table.SetModeTerminal()
	table.AddHeaders("ID", "CreateTime", "Content")
	for _, v := range dd["List"].([]interface{}) {
		va := v.(map[string]interface{})
		table.AddRow(uint(va["ID"].(float64)), va["CreatedAt"], va["Context"])
	}
	fmt.Println(table.Render())
	fmt.Println(dd["TotlePage"], " pages,", dd["TweNum"], "tweets, current ", dd["CurPage"], "page")
}
