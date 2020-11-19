package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add content to the blog",
	Long: `You can add articles, single 
pages, Tweets and categories to the blog`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initAdd() {
	addCmd.PersistentFlags().BoolP("article", "a", false, "Select action article")
	addCmd.PersistentFlags().BoolP("page", "p", false, "Select action page")
	addCmd.PersistentFlags().BoolP("tweet", "t", false, "Select action tweet")
	addCmd.PersistentFlags().BoolP("category", "c", false, "Select action category")
}
