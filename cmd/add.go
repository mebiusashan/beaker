package cmd

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/mebiusashan/beaker/cli"
	"github.com/mebiusashan/beaker/common"
	"github.com/spf13/cobra"
)

var (
	addArticleCatId  uint
	addArticleTitle  string
	addTweetMsg      string
	addCategoryAlias string

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add content to the blog",
		Long: `You can add articles, single 
pages, Tweets and categories to the blog`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	addArticleCmd = &cobra.Command{
		Use:   "article",
		Short: "Add an article",
		Long: `The path of the markdown file 
needs to be set. If the - t parameter is set, 
the parameter value is used as the article title. 
If not set, the file name is used as the article 
title. At the same time, it must be set to the 
ID of the classification to which the chapter belongs`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			if addArticleCatId == 0 {
				common.Err("Need to set category ID")
			}
			mdPath := args[0]
			content, err := ioutil.ReadFile(mdPath)
			common.Assert(err)
			title := addArticleTitle
			if title == "" {
				filenameWithSuffix := path.Base(mdPath)
				fileSuffix := path.Ext(filenameWithSuffix)
				title = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			}
			cli.ArtAdd(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), string(content), title, addArticleCatId)
		},
	}

	addPageCmd = &cobra.Command{
		Use:   "page",
		Short: "Add a page",
		Long: `The path of the markdown file 
needs to be set. If the - t parameter is set, 
the parameter value is used as the article title. 
If not set, the file name is used as the article 
title. At the same time, it must be set to the 
ID of the classification to which the chapter belongs`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()

			mdPath := args[0]
			content, err := ioutil.ReadFile(mdPath)
			common.Assert(err)
			title := addArticleTitle
			if title == "" {
				filenameWithSuffix := path.Base(mdPath)
				fileSuffix := path.Ext(filenameWithSuffix)
				title = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			}
			cli.PageAdd(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), string(content), title)
		},
	}

	addTweetCmd = &cobra.Command{
		Use:   "tweet",
		Short: "Add a tweet",
		Long: `To send a tweet immediately, you 
need to write the content to be sent in the message tag`,
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()

			if addTweetMsg == "" {
				common.Err("Message cannot be empty")
			}
			cli.TweetAdd(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), addTweetMsg)
		},
	}

	addCategoryCmd = &cobra.Command{
		Use:   "category",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()

			if addCategoryAlias == "" {
				common.Err("Alias cannot be empty")
			}
			cli.CatAdd(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), args[0], addCategoryAlias)
		},
	}
)

func init() {
	addArticleCmd.PersistentFlags().UintVarP(&addArticleCatId, "catid", "i", 0, "category ID of the article")
	addArticleCmd.PersistentFlags().StringVarP(&addArticleTitle, "title", "t", "", "title of the article")
	addPageCmd.PersistentFlags().StringVarP(&addArticleTitle, "title", "t", "", "title of the page")
	addTweetCmd.PersistentFlags().StringVarP(&addTweetMsg, "message", "m", "", "message of the tweet")
	addCategoryCmd.PersistentFlags().StringVarP(&addCategoryAlias, "alias", "a", "", "alias of the category")

	addCmd.AddCommand(addArticleCmd)
	addCmd.AddCommand(addPageCmd)
	addCmd.AddCommand(addTweetCmd)
	addCmd.AddCommand(addCategoryCmd)
}
