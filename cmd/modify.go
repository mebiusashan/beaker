package cmd

import (
	"github.com/mebiusashan/beaker/cli"
	"github.com/mebiusashan/beaker/common"
	"github.com/spf13/cobra"
)

var (
	modifyCatId    uint
	modifyCatAlias string

	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "modify content to the blog",
		Long: `You can modify articles, single 
pages, and categories to the blog`,
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
		},
	}

	modifyCatCmd = &cobra.Command{
		Use:   "category",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			if modifyCatId == 0 {
				common.Err("ID cannot be empty")
			}
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			alias := modifyCatAlias
			if name == "" && alias == "" {
				common.Err("There are no changes")
			}
			cli.CatModify(getWebsiteInfo().HOST, modifyCatId, name, alias)
		},
	}
)

func init() {
	addArticleCmd.PersistentFlags().UintVarP(&modifyCatId, "catid", "i", 0, "category ID of the article")
	addCategoryCmd.PersistentFlags().StringVarP(&modifyCatAlias, "alias", "a", "", "alias of the category")
}
