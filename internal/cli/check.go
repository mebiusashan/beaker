package cli

import (
	"strings"

	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/net"
)

func Check(host string, serverDesKey []byte) bool {
	T := []byte("beaker------====------")
	desT, err := cert.TripleDesEncrypt(T, serverDesKey)
	jsonData := net.PostJson(host+net.CLI_CHECK, strings.NewReader(cert.Base64Encode(desT)))
	rel64, err := cert.Base64Decode(jsonData.Data.(string))
	common.Assert(err)
	checkT, err := cert.TripleDesDecrypt(rel64, serverDesKey)
	common.Assert(err)
	if string(checkT) == (string(T) + "123") {
		return true
	}
	return false
}
