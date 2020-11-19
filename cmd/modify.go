package cmd

import (
	"github.com/spf13/cobra"
)

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "modify content to the blog",
	Long: `You can modify articles, single 
pages, and categories to the blog`,
	Run: func(cmd *cobra.Command, args []string) {
		checkWebsite()
	},
}

func init() {
	modifyCmd.PersistentFlags().BoolP("article", "a", false, "Select action article")
	modifyCmd.PersistentFlags().BoolP("page", "p", false, "Select action page")
	modifyCmd.PersistentFlags().BoolP("tweet", "t", false, "Select action tweet")
	modifyCmd.PersistentFlags().BoolP("category", "c", false, "Select action category")
}
