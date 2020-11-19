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

func PageAll(host string) {
	resp, err := http.Post(host+"/admin/pag/list", "", strings.NewReader(""))
	common.Assert(err)

	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
	}

	table := termtables.CreateTable()
	table.AddHeaders("ID", "Title")
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		table.AddRow(uint(v["ID"].(float64)), v["Title"])
	}
	fmt.Println(table.Render())
}
