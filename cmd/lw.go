package cmd

import (
	"github.com/spf13/cobra"
)

var lwCmd = &cobra.Command{
	Use:   "lw",
	Short: "View blog list",
	Long: `This command can view all 
blog information stored locally`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
