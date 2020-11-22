package cli

import (
	"fmt"
	"strconv"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func CatAll(host string, refresh bool, key []byte) {
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_CAT_LIST, refresh, key, "")

	maxid := 0
	maxcname := 0
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		l := len(strconv.Itoa(int(v["ID"].(float64))))
		if maxid < l {
			maxid = l
		}
		if maxcname < len(v["Alias"].(string)) {
			maxcname = len(v["Alias"].(string))
		}
	}
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-"+strconv.Itoa(maxid)+"d %-"+strconv.Itoa(maxcname)+"s %s\n", uint(v["ID"].(float64)), v["Alias"], v["Name"])
	}
}

func CatRm(host string, refresh bool, key []byte, id uint, mid uint) {
	sendData := common.CatRmReq{}
	sendData.ID = id
	sendData.MvID = mid
	net.PostJsonWithEncrypt(host+net.CLI_CAT_RM, refresh, key, sendData)
}

func CatAdd(host string, refresh bool, key []byte, name string, alias string) {
	sendData := common.CatModel{Alias: alias, Name: name}
	net.PostJsonWithEncrypt(host+net.CLI_CAT_ADD, refresh, key, sendData)
}

func CatModify(host string, refresh bool, key []byte, id uint, name string, alias string) {
	sendData := common.CatModel{Alias: alias, Name: name}
	sendData.ID = id
	net.PostJsonWithEncrypt(host+net.CLI_CAT_MODIFY, refresh, key, sendData)
}
