package cmd

import (
	"github.com/spf13/cobra"
)

var (
	actionWebsite string

	rootCmd = &cobra.Command{
		Use:   "beaker",
		Short: "Beaker is a simple blog system.",
		Long: `Beaker is a CS architecture blog system, 
you can manage your numerous beaker blogs through beaker.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&actionWebsite, "website", "w", "default", "Set the blog you want to push, the blog can be set in the config command")
	rootCmd.PersistentFlags().BoolP("refresh", "r", true, "refresh server cache")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(modifyCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(lwCmd)
}
