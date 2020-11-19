package cmd

import (
	"strconv"

	"github.com/mebiusashan/beaker/cli"
	"github.com/mebiusashan/beaker/common"
	"github.com/spf13/cobra"
)

var (
	articlerm  bool
	pagerm     bool
	tweetrm    bool
	categoryrm bool

	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Delete some content from the blog",
		Long:  `You can delete articles, single pages, Tweets and categories in the blog`,
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			if articlerm {
				if len(args) == 0 {
					common.Err("Missing ID parameter")
				}
				id, err := strconv.Atoi(args[0])
				common.Assert(err)
				cli.ArtRm(getWebsiteInfo().HOST, uint(id))
				return
			}

		},
	}
)

func init() {
	rmCmd.PersistentFlags().BoolVarP(&articlerm, "article", "a", false, "Select action article")
	rmCmd.PersistentFlags().BoolVarP(&pagerm, "page", "p", false, "Select action page")
	rmCmd.PersistentFlags().BoolVarP(&tweetrm, "tweet", "t", false, "Select action tweet")
	rmCmd.PersistentFlags().BoolVarP(&categoryrm, "category", "c", false, "Select action category")
}
