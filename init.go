package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use: "init",
	Run: func(command *cobra.Command, args []string) {
		file, err := os.Create(versionFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		tag, err := getHeadTag("")
		switch err {
		case nil:
			_, err := file.WriteString(tag.String())
			if err != nil {
				fmt.Println(err)
				return
			}
		case ErrNoSemverTags:
			version := semver.MustParse("0.0.1")
			_, err := file.WriteString(version.String())
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Println(err)
			return
		}
		return
	},
}
