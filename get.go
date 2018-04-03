package main

import (
	"fmt"
	"strings"

	"github.com/ninedraft/semtag/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCurrentVersionAliases = []string{"tag", "tags"}
var getCurrentVersionConfig struct {
	All bool
}
var getCurrentVersion = &cobra.Command{
	Use:     "show",
	Aliases: getCurrentVersionAliases,
	PreRun: func(command *cobra.Command, args []string) {

	},
	Run: func(command *cobra.Command, args []string) {
		if getCurrentVersionConfig.All {
			repo := git.Repo{}
			tags, err := repo.Tags()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(strings.Join(tags, "\n"))
			return
		}
		version, err := getHeadTag("")
		switch {
		case err == nil:
			// pass
		default:
			fmt.Println(err)
			return
		}
		fmt.Printf("v%v\n", version)
	},
}

func init() {
	getCurrentVersion.PersistentFlags().BoolVarP(&getCurrentVersionConfig.All, "all", "a", false, "show all versions")
	viper.BindPFlag("all", increment.PersistentFlags().Lookup("all"))
}
