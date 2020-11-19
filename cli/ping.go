package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mebiusashan/beaker/cert"
	"github.com/mebiusashan/beaker/common"
)

func Ping(host string) []byte {
	resp, err := http.Post(host+"/user/ping", "", strings.NewReader(""))
	common.Assert(err)

	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
	}

	pubKeyStr := jsonData.Data.(string)
	pubkey, err := cert.Base64Decode(pubKeyStr)
	common.Assert(err)
	return pubkey
}
