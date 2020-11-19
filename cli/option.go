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

func OptAll(host string) {
	resp, err := http.Post(host+"/admin/opt", "", strings.NewReader(""))
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

	table := termtables.CreateTable()
	table.AddHeaders("name", "value")
	for k, v := range jsonData.Data.(map[string]interface{}) {
		table.AddRow(k, v)
	}
	fmt.Println(table.Render())
}
