package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func ArtAll(host string) {
	jsonData := net.PostJson(host+net.CLI_ART_LIST, strings.NewReader(""))

	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		fmt.Printf("%-5d%s\n", uint(v["ID"].(float64)), v["Title"])
	}
}

func ArtRm(host string, id uint) {
	sendData := common.ArcDB{}
	sendData.ID = id
	jsonByte, err := json.Marshal(sendData)
	common.Assert(err)

	net.PostJson(host+net.CLI_ART_RM, strings.NewReader(string(jsonByte)))
}
