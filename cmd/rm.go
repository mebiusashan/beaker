package cmd

import (
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete some content from the blog",
	Long:  `You can delete articles, single pages, Tweets and categories in the blog`,
	Run: func(cmd *cobra.Command, args []string) {
		checkWebsite()
	},
}

func init() {
	rmCmd.PersistentFlags().BoolP("article", "a", false, "Select action article")
	rmCmd.PersistentFlags().BoolP("page", "p", false, "Select action page")
	rmCmd.PersistentFlags().BoolP("tweet", "t", false, "Select action tweet")
	rmCmd.PersistentFlags().BoolP("category", "c", false, "Select action category")
}
