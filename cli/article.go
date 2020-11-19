package cli

import (
	"fmt"
	"strings"

	"github.com/mebiusashan/beaker/net"
)

func ArtAll(host string) {
	jsonData := net.PostJson(host+"/admin/arc/list", strings.NewReader(""))

	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-5d%s\n", uint(v["ID"].(float64)), v["Title"])
	}
}
