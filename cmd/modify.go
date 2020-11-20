package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/mebiusashan/beaker/cli"
	"github.com/mebiusashan/beaker/common"
	"github.com/spf13/cobra"
)

var (
	modifyCatId    uint
	modifyCatAlias string
	modifyTitle    string

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
		Short: "asdf",
		Long:  `asdfa`,
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

	modifyArtCmd = &cobra.Command{
		Use:   "article",
		Short: "sda",
		Long:  `afas`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// checkWebsite()
			id, err := strconv.Atoi(args[0])
			common.Assert(err)
			curPath, err := os.Getwd()
			common.Assert(err)
			title, content := cli.ArtDownload(getWebsiteInfo().HOST, uint(id))
			curPath = curPath + "/" + title + ".md"
			err = ioutil.WriteFile(curPath, []byte(content), 0666)
			common.Assert(err)
			runEditor(curPath)
			//edit file complete, push file
			newContent, err := ioutil.ReadFile(curPath)
			common.Assert(err)
			cli.ArtModify(getWebsiteInfo().HOST, uint(id), modifyCatId, modifyTitle, string(newContent))
		},
	}
)

func init() {

	modifyCatCmd.PersistentFlags().UintVarP(&modifyCatId, "catid", "i", 0, "category ID of the category")
	modifyCatCmd.PersistentFlags().StringVarP(&modifyCatAlias, "alias", "a", "", "alias of the category")

	modifyArtCmd.PersistentFlags().UintVarP(&modifyCatId, "catid", "i", 0, "category ID of the article")
	modifyArtCmd.PersistentFlags().StringVarP(&modifyTitle, "title", "t", 0, "title of the article")

	modifyCmd.AddCommand(modifyArtCmd)
	modifyCmd.AddCommand(modifyCatCmd)
}

func runEditor(path string) {
	editor := common.DefaultEditor
	if localConfig.Editor != "" {
		editor = localConfig.Editor
	}
	oscmd := exec.Command(editor, path)
	oscmd.Stdin = os.Stdin
	oscmd.Stdout = os.Stdout
	err := oscmd.Run()
	common.Assert(err)
}
