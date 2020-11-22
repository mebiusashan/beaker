package cmd

import (
	"fmt"
	"strconv"

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
	Editor         string
	Websites       []websiteConfig
}

func (c *websiteConfig) GetKey() []byte {
	return []byte(c.Key)
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
			editor := common.DefaultEditor
			if localConfig.Editor != "" {
				editor = localConfig.Editor
			}
			fmt.Println("Editor:", editor)
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
			login(url, addWebSiteUser, addWebSitePassword)
		},
	}

	rmWebsiteCmd = &cobra.Command{
		Use:   "rmw",
		Short: "Remove a website information",
		Long: `Remove a website information in 
the local record, if the website is the default 
website, set the first one of the remaining 
website information as the default website`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			alias := args[0]
			isDefault := localConfig.DefaultWebsite == alias
			for i := 0; i < len(localConfig.Websites); i++ {
				if isDefault && localConfig.Websites[i].Alias != alias {
					localConfig.DefaultWebsite = localConfig.Websites[i].Alias
					isDefault = false
				}
				if localConfig.Websites[i].Alias == alias {
					localConfig.Websites = append(localConfig.Websites[:i], localConfig.Websites[i+1:]...)
					writeConfig()
					return
				}
			}
		},
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to the website",
		Long: `When the website login expires, 
use the login command to log in again and obtain 
operation permissions`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			alias := args[0]
			info := getWebsiteInfoWithAlias(alias)
			login(info.HOST, addWebSiteUser, addWebSitePassword)
		},
	}

	setEditorCmd = &cobra.Command{
		Use:   "editor",
		Short: "Set up a text editor",
		Long: `When modifying an article, the content 
of the article is downloaded to the local, and then 
opened and edited with this editor`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			editor := args[0]
			localConfig.Editor = editor
			writeConfig()
		},
	}
)

func init() {
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSiteAlias, "alias", "a", "", "blog alias")
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSiteUser, "user", "u", "", "blog administrator account name")
	addWebSiteCmd.PersistentFlags().StringVarP(&addWebSitePassword, "password", "p", "", "blog administrator account password")
	addWebSiteCmd.PersistentFlags().BoolVarP(&addWebSiteIsDefault, "defalut", "d", false, "set as default blog")

	loginCmd.PersistentFlags().StringVarP(&addWebSiteUser, "user", "u", "", "blog administrator account name")
	loginCmd.PersistentFlags().StringVarP(&addWebSitePassword, "password", "p", "", "blog administrator account password")

	configCmd.AddCommand(addWebSiteCmd)
	configCmd.AddCommand(rmWebsiteCmd)
	configCmd.AddCommand(loginCmd)
	configCmd.AddCommand(setEditorCmd)
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

func checkWebsite() {
	if len(localConfig.Websites) == 0 {
		common.Err("No website information is configured")
	}
	info := getWebsiteInfo()
	rel := cli.Check(info.HOST, []byte(info.Key))
	if !rel {
		fmt.Println("Login invalid, please login again")
	}
}

func getWebsiteInfo() *websiteConfig {
	return getWebsiteInfoWithAlias(actionWebsite)
}

func getWebsiteInfoWithAlias(alias string) *websiteConfig {
	if alias == "" || alias == "default" {
		alias = localConfig.DefaultWebsite
	}

	for _, d := range localConfig.Websites {
		if d.Alias == alias {
			return &d
		}
	}

	common.Err(actionWebsite + " website information does not exist")
	return nil
}

func login(url string, username string, password string) {
	//need login with current website
	//if success login, set info and server key to config file
	//if it's default website, need reset all website config
	pubKey := cli.Ping(url)
	serverPubKey := cli.Login(url, pubKey, username, password)
	if cli.Check(url, serverPubKey) {
		if len(localConfig.Websites) == 0 || localConfig.Websites == nil {
			localConfig.Websites = make([]websiteConfig, 0)
			addWebSiteIsDefault = true
		}
		if addWebSiteIsDefault {
			localConfig.DefaultWebsite = addWebSiteAlias
		}

		for i := 0; i < len(localConfig.Websites); i++ {
			v := localConfig.Websites[i]
			if v.HOST == url {
				localConfig.Websites[i].Key = string(serverPubKey)
				writeConfig()
				return
			}
		}

		d := websiteConfig{Alias: addWebSiteAlias, HOST: url, Key: string(serverPubKey)}
		localConfig.Websites = append(localConfig.Websites, d)
		writeConfig()
	} else {
		fmt.Println("login fail")
	}
}

func writeConfig() {
	viper.Set("config", localConfig)
	err := viper.WriteConfig()
	common.Assert(err)
}
