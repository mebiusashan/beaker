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

	max := 0
	for k := range jsonData.Data.(map[string]interface{}) {
		if len(k) > max {
			max = len(k)
		}
	}
	for k, v := range jsonData.Data.(map[string]interface{}) {
		fmt.Printf("%-"+strconv.Itoa(max)+"s ", k)
		fmt.Println(v)
	}
}
