package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/net"
)

func CatAll(host string) {
	jsonData := net.PostJson(host+"/admin/cat/list", strings.NewReader(""))

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
