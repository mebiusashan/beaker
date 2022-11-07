package command

import (
	"fmt"

	"github.com/mebiusashan/beaker/internal/common"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Beaker",
	Long:  `All software has versions. This is Beaker's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(common.VERSION)
	},
}
