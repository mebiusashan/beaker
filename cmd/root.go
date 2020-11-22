package cmd

import (
	"github.com/spf13/cobra"
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
code can be accessed at https://github.com/mebiusashan/beaker.`,
	}
)

// Execute executes the root command.
func Execute() error {
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
}
