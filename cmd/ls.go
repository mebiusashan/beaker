package cmd

import (
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "View server data list",
	Long: `You can view the corresponding data 
list according to the flag, and the server will 
return all the corresponding data`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initLs() {
	lsCmd.PersistentFlags().BoolP("article", "a", false, "Select action article")
	lsCmd.PersistentFlags().BoolP("page", "p", false, "Select action page")
	lsCmd.PersistentFlags().BoolP("tweet", "t", false, "Select action tweet")
	lsCmd.PersistentFlags().BoolP("category", "c", false, "Select action category")
	lsCmd.PersistentFlags().BoolP("option", "o", false, "Select action option")
}
