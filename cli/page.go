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
