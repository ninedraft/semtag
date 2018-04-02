package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
		_, err := os.Stat(versionFileName)
		switch {
		case os.IsNotExist(err):
			fmt.Printf("Can't update version in unitialized project!\nRun semtag init")
			return
		case err == nil:
		default:
			fmt.Println(err)
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
			err := ioutil.WriteFile(versionFileName, []byte(version.String()), os.ModePerm)
			if err != nil {
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
