package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear server cache.",
	Long: `Every time a user visits a blog, 
the visited pages will be stored in the 
server's cache. If you modify the blog's 
template or some settings, you need to 
execute clean to clear the server's cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}
