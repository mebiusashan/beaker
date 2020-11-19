package cmd

import (
	"github.com/mebiusashan/beaker/cli"
	"github.com/spf13/cobra"
)

var (
	articlels  bool
	pagels     bool
	tweetls    bool
	categoryls bool
	optionls   bool

	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "View server data list",
		Long: `You can view the corresponding data 
list according to the flag, and the server will 
return all the corresponding data`,
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			if articlels {
				cli.ArtAll(getWebsiteInfo(actionWebsite).HOST)
			}
			if pagels {
				cli.PageAll(getWebsiteInfo(actionWebsite).HOST)
			}
			if tweetls {
				cli.TweetAll(getWebsiteInfo(actionWebsite).HOST)
			}
			if categoryls {
				cli.CatAll(getWebsiteInfo(actionWebsite).HOST)
			}
			if optionls {
				cli.OptAll(getWebsiteInfo(actionWebsite).HOST)
			}
		},
	}
)

func init() {
	lsCmd.PersistentFlags().BoolVarP(&articlels, "article", "a", false, "Select action article")
	lsCmd.PersistentFlags().BoolVarP(&pagels, "page", "p", false, "Select action page")
	lsCmd.PersistentFlags().BoolVarP(&tweetls, "tweet", "t", false, "Select action tweet")
	lsCmd.PersistentFlags().BoolVarP(&categoryls, "category", "c", false, "Select action category")
	lsCmd.PersistentFlags().BoolVarP(&optionls, "option", "o", false, "Select action option")
}
