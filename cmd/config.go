package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure blog information",
	Long: `The configuration command can 
set your blog background address`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
