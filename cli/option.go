package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebiusashan/beaker/net"
)

func OptAll(host string) {
	jsonData := net.PostJson(host+net.CLI_OPTION, strings.NewReader(""))

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
