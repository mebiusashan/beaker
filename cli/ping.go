package cli

import (
	"strings"

	"github.com/mebiusashan/beaker/cert"
	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func Ping(host string) []byte {
	jsonData := net.PostJson(host+"/user/ping", strings.NewReader(""))

	pubKeyStr := jsonData.Data.(string)
	pubkey, err := cert.Base64Decode(pubKeyStr)
	common.Assert(err)
	return pubkey
}
