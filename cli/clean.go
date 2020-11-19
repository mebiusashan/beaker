package cli

import (
	"strings"

	"github.com/mebiusashan/beaker/net"
)

func CleanCache(host string) {
	net.PostJson(host+net.CLI_CLEAN, strings.NewReader(""))
}
