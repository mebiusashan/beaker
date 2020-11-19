package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var lwCmd = &cobra.Command{
	Use:   "lw",
	Short: "View blog list",
	Long: `This command can view all 
blog information stored locally`,
	Run: func(cmd *cobra.Command, args []string) {
		l := len(localConfig.Websites)
		fmt.Println("Total number of sites:", l)
		if l == 0 {
			return
		}
		fmt.Println("Default Website:", localConfig.DefaultWebsite)
		fmt.Println("-----------------------------")
		max := 0
		for _, d := range localConfig.Websites {
			if len(d.Alias) > max {
				max = len(d.Alias)
			}
		}
		for _, d := range localConfig.Websites {
			fmt.Printf("%-"+strconv.Itoa(max)+"s %s\n", d.Alias, d.HOST)
		}
	},
}
