package cli

import (
	"github.com/mebiusashan/beaker/net"
)

func CleanCache(host string, refresh bool, key []byte) {
	net.PostJsonWithEncrypt(host+net.CLI_CLEAN, refresh, key, "")
}
