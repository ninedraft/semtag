package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCurrentVersionAliases = []string{}
var getCurrentVersionConfig struct {
}
var getCurrentVersion = &cobra.Command{
	Use:     "get",
	Aliases: getCurrentVersionAliases,
	PreRun: func(command *cobra.Command, args []string) {

	},
	Run: func(command *cobra.Command, args []string) {
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
