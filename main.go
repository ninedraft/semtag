package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

const (
	versionFileName = "version.txt"
)

var root = &cobra.Command{
	Use: "semtag",
	Run: func(command *cobra.Command, args []string) {

	},
}

func main() {
	root.AddCommand(
		increment,
		initCommand,
	)
	err := root.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func readSemverFile() (semver.Version, error) {
	data, err := ioutil.ReadFile(versionFileName)
	if err != nil {
		return semver.Version{}, err
	}
	semverString := strings.TrimSpace(string(data))
	version, err := semver.Parse(semverString)
	return version, err
}