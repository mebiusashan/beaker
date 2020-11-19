package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mebiusashan/beaker/common"

	"github.com/mebiusashan/beaker/cert"
)

func Login(host string, pubKey []byte, username string, password string) []byte {
	clientDesKey := cert.CreateDesKey()
	clientDesKeyM, err := cert.RSAEncryp(pubKey, clientDesKey)
	common.Assert(err)

	jsonSendData := common.LoginReq{DK: cert.Base64Encode(clientDesKeyM), UN: username, PW: cert.MD5([]byte(password))}
	jsonByte, err := json.Marshal(jsonSendData)
	common.Assert(err)

	jsonStr := cert.Base64Encode(jsonByte)
	resp, err := http.Post(host+"/user/login", "", strings.NewReader(jsonStr))
	common.Assert(err)

	body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
	}

	serverDesKeyM, err := cert.Base64Decode(jsonData.Data.(string))
	common.Assert(err)
	serverDesKey, err := cert.TripleDesDecrypt(serverDesKeyM, clientDesKey)
	common.Assert(err)

	return serverDesKey
}
