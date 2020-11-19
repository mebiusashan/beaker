package cmd

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/mebiusashan/beaker/cli"
	"github.com/mebiusashan/beaker/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type websiteConfig struct {
	Alias string
	HOST  string
	Key   string
}

type config struct {
	DefaultWebsite string
	Websites       []websiteConfig
}

var (
	localConfig         config
	addWebSiteAlias     string
	addWebSiteUser      string
	addWebSitePassword  string
	addWebSiteIsDefault bool

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configure blog information",
		Long: `The configuration command can 
set your blog background address`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	addWebSiteCmd = &cobra.Command{
		Use:   "addw",
		Short: "Add blog site information",
		Long: `To add a blog server information, 
you need to specify the blog address, 
user name and password, and it will be 
verified after adding`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			if !govalidator.IsURL(url) && !govalidator.IsIP(url) {
				common.Err("Blog address format error")
			}

			if addWebSiteAlias == "" {
				common.Err("Alias ​​cannot be empty")
			}

			//check current config , has duplicate config
			for _, website := range localConfig.Websites {
				if website.Alias == addWebSiteAlias {
					common.Err("Duplicate alias")
				}
				if website.HOST == url {
					common.Err("Duplicate host")
				}
			}

			//need login with current website
			//if success login, set info and server key to config file
			//if it's default website, need reset all website config
			pubKey := cli.Ping(url)
			serverPubKey := cli.Login(url, pubKey, addWebSiteUser, addWebSitePassword)
			if cli.Check(url, serverPubKey) {
				if len(localConfig.Websites) == 0 || localConfig.Websites == nil {
					localConfig.Websites = make([]websiteConfig, 0)
					addWebSiteIsDefault = true
				}
				if addWebSiteIsDefault {
					localConfig.DefaultWebsite = addWebSiteAlias
				}
				d := websiteConfig{Alias: addWebSiteAlias, HOST: url, Key: string(serverPubKey)}

				localConfig.Websites = append(localConfig.Websites, d)
				viper.Set("config", localConfig)
				err := viper.WriteConfig()
				common.Assert(err)
			} else {
				fmt.Println("login fail")
			}
		},
	}
)

func init() {
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSiteAlias, "alias", "a", "", "blog alias")
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSiteUser, "user", "u", "", "blog administrator account name")
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSitePassword, "password", "p", "", "blog administrator account password")
	addWebSiteCmd.PersistentFlags().BoolVarP(&addWebSiteIsDefault, "defalut", "d", false, "set as default blog")

	configCmd.AddCommand(addWebSiteCmd)
}

func initConfig() {
	home, err := homedir.Dir()
	common.Assert(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".beaker")
	viper.SafeWriteConfig()

	err = viper.ReadInConfig()
	common.Assert(err)

	err = viper.UnmarshalKey("config", &localConfig)
	common.Assert(err)
}
