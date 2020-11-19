package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func CatAll(host string) {
	jsonData := net.PostJson(host+net.CLI_CAT_LIST, strings.NewReader(""))

	maxid := 0
	maxcname := 0
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		l := len(strconv.Itoa(int(v["ID"].(float64))))
		if maxid < l {
			maxid = l
		}
		if maxcname < len(v["Cname"].(string)) {
			maxcname = len(v["Cname"].(string))
		}
	}
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-"+strconv.Itoa(maxid)+"d %-"+strconv.Itoa(maxcname)+"s %s\n", uint(v["ID"].(float64)), v["Cname"], v["Name"])
	}
}

func CatRm(host string, id uint, mid uint) {
	sendData := common.CatDBDel{}
	sendData.ID = id
	sendData.MvID = mid
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_CAT_RM, strings.NewReader(string(jsonByte)))
}

func CatAdd(host string, name string, alias string) {

}
