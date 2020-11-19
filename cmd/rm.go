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
			if pagerm {
				if len(args) == 0 {
					common.Err("Missing ID parameter")
				}
				id, err := strconv.Atoi(args[0])
				common.Assert(err)
				cli.PageRm(getWebsiteInfo().HOST, uint(id))
				return
			}
			if tweetrm {
				if len(args) == 0 {
					common.Err("Missing ID parameter")
				}
				id, err := strconv.Atoi(args[0])
				common.Assert(err)
				cli.TweetRm(getWebsiteInfo().HOST, uint(id))
				return
			}

			if categoryrm {
				if len(args) < 2 {
					common.Err("Missing ID parameter")
				}
				id, err := strconv.Atoi(args[0])
				common.Assert(err)
				mid, err := strconv.Atoi(args[1])
				common.Assert(err)
				cli.CatRm(getWebsiteInfo().HOST, uint(id), uint(mid))
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
