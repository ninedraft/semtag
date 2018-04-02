package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ninedraft/semtag/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var incrementAliases = []string{"inc", "incr", "+"}
var incrementConfig struct {
	Force bool
}
var increment = &cobra.Command{
	Use:     "increment",
	Aliases: incrementAliases,
	PreRun: func(command *cobra.Command, args []string) {

	},
	Run: func(command *cobra.Command, args []string) {
		if len(args) == 0 {
			command.Help()
			return
		}
		level := args[0]
		version, err := getHeadTag("")
		switch {
		case err == nil:
			// pass
		case os.IsNotExist(err):
			return
		default:
			fmt.Println(err)
			return
		}
		switch level {
		case "patch":
			version.Patch++
		case "minor":
			version.Minor++
		case "major":
			version.Major++
		default:
			command.Help()
			return
		}
		if incrementConfig.Force || areYouSure() {
			if err := storeLocalVersion(version); err != nil {
				fmt.Println(err)
				return
			}
			if err := (&git.Repo{}).AddTag(version); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(version)
		}
	},
}

func init() {
	increment.PersistentFlags().BoolVarP(&incrementConfig.Force, "force", "f", false, "don't ask confirmation")
	viper.BindPFlag("force", increment.PersistentFlags().Lookup("force"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "MIT")
}

func areYouSure() bool {
	fmt.Printf("Are you sure?[Y/N]:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			panic(err)
		}
		break
	}
	return strings.ToLower(strings.TrimSpace(scanner.Text())) == "y"
}
