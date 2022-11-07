package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/mebiusashan/beaker/internal/cli"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/spf13/cobra"
)

var (
	modifyCatId         uint
	modifyCatAlias      string
	modifyTitle         string
	modifyDelFile       bool
	modifyIgnoreContent bool

	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "modify content to the blog",
		Long: `You can modify articles, single 
pages, and categories to the blog`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	modifyCatCmd = &cobra.Command{
		Use:   "category",
		Short: "Modify category information",
		Long: `To modify the existing classification information, 
you need to set the ID of the classification, and then you can 
modify the name and alias`,
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
			cli.CatModify(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), modifyCatId, name, alias)
		},
	}

	modifyArtCmd = &cobra.Command{
		Use:   "article",
		Short: "modify articles",
		Long: `According to the entered article ID, the article 
will be downloaded to the local first, and then opened with the 
set text editor. When the user exits the text editor after the 
modification is completed, the program will automatically synchronize 
the modified text content to the server, and you can also modify the 
article title and category`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			id, err := strconv.Atoi(args[0])
			common.Assert(err)
			mdStr := ""
			curPath := ""
			var imgs []common.ImgInfo
			if !modifyIgnoreContent {
				curPath, err := os.Getwd()
				common.Assert(err)
				title, content := cli.ArtDownload(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), uint(id))
				curPath = curPath + "/" + title + ".md"
				err = ioutil.WriteFile(curPath, []byte(content), 0666)
				common.Assert(err)
				runEditor(curPath)
				//edit file complete, push file
				newContent, err := ioutil.ReadFile(curPath)
				common.Assert(err)
				mdStr, imgs = convMarkdownImage(newContent, curPath)
			}
			cli.ArtModify(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), uint(id), modifyCatId, modifyTitle, mdStr, imgs)
			if !modifyIgnoreContent && modifyDelFile {
				err = os.Remove(curPath)
				common.Assert(err)
			}
		},
	}

	modifyPageCmd = &cobra.Command{
		Use:   "page",
		Short: "modify page",
		Long: `According to the entered page ID, the page 
will be downloaded to the local first, and then opened with the 
set text editor. When the user exits the text editor after the 
modification is completed, the program will automatically synchronize 
the modified text content to the server, and you can also modify the 
page title`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkWebsite()
			id, err := strconv.Atoi(args[0])
			common.Assert(err)
			mdStr := ""
			curPath := ""
			var imgs []common.ImgInfo
			if !modifyIgnoreContent {
				curPath, err := os.Getwd()
				common.Assert(err)
				title, content := cli.PageDownload(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), uint(id))
				curPath = curPath + "/" + title + ".md"
				err = ioutil.WriteFile(curPath, []byte(content), 0666)
				common.Assert(err)
				runEditor(curPath)
				//edit file complete, push file
				newContent, err := ioutil.ReadFile(curPath)
				common.Assert(err)
				mdStr, imgs = convMarkdownImage(newContent, curPath)
			}
			cli.PageModify(getWebsiteInfo().HOST, refresh, getWebsiteInfo().GetKey(), uint(id), modifyTitle, mdStr, imgs)
			if !modifyIgnoreContent && modifyDelFile {
				err = os.Remove(curPath)
				common.Assert(err)
			}
		},
	}
)

func init() {

	modifyCatCmd.PersistentFlags().UintVarP(&modifyCatId, "catid", "i", 0, "category ID of the category")
	modifyCatCmd.PersistentFlags().StringVarP(&modifyCatAlias, "alias", "a", "", "alias of the category")
	modifyCatCmd.MarkFlagRequired("catid")

	modifyArtCmd.PersistentFlags().UintVarP(&modifyCatId, "catid", "i", 0, "category ID of the article")
	modifyArtCmd.PersistentFlags().StringVarP(&modifyTitle, "title", "t", "", "title of the article")
	modifyArtCmd.PersistentFlags().BoolVarP(&modifyDelFile, "del", "d", false, "Delete locally cached markdown files")
	modifyArtCmd.PersistentFlags().BoolVarP(&modifyIgnoreContent, "ignore", "g", false, "Ignore content")

	modifyPageCmd.PersistentFlags().StringVarP(&modifyTitle, "title", "t", "", "title of the page")
	modifyPageCmd.PersistentFlags().BoolVarP(&modifyDelFile, "del", "d", false, "Delete locally cached markdown files")
	modifyPageCmd.PersistentFlags().BoolVarP(&modifyIgnoreContent, "ignore", "g", false, "Ignore content")

	modifyCmd.AddCommand(modifyArtCmd)
	modifyCmd.AddCommand(modifyCatCmd)
	modifyCmd.AddCommand(modifyPageCmd)
}

func runEditor(path string) {
	editor := common.DefaultEditor
	if localConfig.Editor != "" {
		editor = localConfig.Editor
	}
	fmt.Println(editor, path)
	oscmd := exec.Command(editor, path)
	oscmd.Stdin = os.Stdin
	oscmd.Stdout = os.Stdout
	err := oscmd.Run()
	fmt.Println(err)
	common.Assert(err)
}
