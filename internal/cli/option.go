package cli

import (
	"fmt"
	"strconv"

	"github.com/mebiusashan/beaker/internal/net"
)

func OptAll(host string, refresh bool, key []byte) {
	jsonData := net.PostJsonWithEncrypt(host+net.CLI_OPTION, refresh, key, "")

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
