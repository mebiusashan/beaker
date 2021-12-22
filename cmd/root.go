package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	actionWebsite string
	refresh       bool

	rootCmd = &cobra.Command{
		Use:   "beaker",
		Short: "Beaker is a simple blog system.",
		Long: `Beaker is a very fast, simple and 
smart blog system. It is very suitable for geeks. 
At the same time, Beaker advocates using markdown 
to edit articles and manage your blog through the 
terminal. It is completely open source, you can 
add and modify functions at will, and its source 
code can be accessed at https://beaker.ashan.org.`,
	}
)

// Execute executes the root command.
func Execute() error {
	rootCmd.DisableAutoGenTag = true
	sysType := runtime.GOOS
	if sysType != "windows" {
		header := &doc.GenManHeader{
			Title:   "Beaker",
			Section: "1",
		}
		doc.GenManTree(rootCmd, header, "/usr/local/share/man/man1")
	}
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&actionWebsite, "website", "w", "default", "Set the blog you want to push, the blog can be set in the config command")
	rootCmd.PersistentFlags().BoolVarP(&refresh, "refresh", "r", true, "refresh server cache")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(modifyCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(lwCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(loginCmd)
}
