package cli

import (
	"encoding/json"
	"strings"

	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/net"

	"github.com/mebiusashan/beaker/internal/cert"
)

func Login(host string, pubKey []byte, username string, password string) []byte {
	clientDesKey := cert.CreateDesKey()
	clientDesKeyM, err := cert.RSAEncryp(pubKey, clientDesKey)
	common.Assert(err)

	jsonSendData := common.LoginReq{DK: cert.Base64Encode(clientDesKeyM), UN: username, PW: cert.MD5([]byte(password))}
	jsonByte, err := json.Marshal(jsonSendData)
	common.Assert(err)

	jsonStr := cert.Base64Encode(jsonByte)
	jsonData := net.PostJson(host+net.CLI_LOGIN, strings.NewReader(jsonStr))

	serverDesKeyM, err := cert.Base64Decode(jsonData.Data.(string))
	common.Assert(err)
	serverDesKey, err := cert.TripleDesDecrypt(serverDesKeyM, clientDesKey)
	common.Assert(err)

	return serverDesKey
}
