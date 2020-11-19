package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mebiusashan/beaker/cert"
	"github.com/mebiusashan/beaker/common"
)

func Check(host string, serverDesKey []byte) bool {
	T := []byte("beaker------====------")
	desT, err := cert.TripleDesEncrypt(T, serverDesKey)
	common.Assert(err)

	resp, err := http.Post(host+"/user/check", "", strings.NewReader(cert.Base64Encode(desT)))
	common.Assert(err)

	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		return false
	}

	rel64, err := cert.Base64Decode(jsonData.Data.(string))
	common.Assert(err)

	checkT, err := cert.TripleDesDecrypt(rel64, serverDesKey)
	common.Assert(err)

	if string(checkT) == (string(T) + "123") {
		return true
	}
	return false
}
