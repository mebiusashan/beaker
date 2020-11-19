package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apcera/termtables"
	"github.com/mebiusashan/beaker/common"
)

func CatAll(host string) {
	resp, err := http.Post(host+"/admin/cat/list", "", strings.NewReader(""))
	common.Assert(err)
	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		fmt.Println(jsonData.Msg)
		return
	}

	table := termtables.CreateTable()
	table.AddHeaders("ID", "showname", "path", "CreateTime")
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		table.AddRow(uint(v["ID"].(float64)), v["Name"], v["Cname"], v["CreatedAt"])
	}
	fmt.Println(table.Render())
}
