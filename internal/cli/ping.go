package cli

import (
	"strings"

	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/net"
)

func Ping(host string) []byte {
	jsonData := net.PostJson(host+net.CLI_PING, strings.NewReader(""))

	pubKeyStr := jsonData.Data.(string)
	pubkey, err := cert.Base64Decode(pubKeyStr)
	common.Assert(err)
	return pubkey
}
