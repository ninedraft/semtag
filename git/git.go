package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/blang/semver"
)

type Repo struct {
	path string
}

func (repo *Repo) FetchTags(repos ...string) error {
	args := []string{"fetch", "--tags"}
	switch len(repos) {
	case 0:
		// pass
	case 1:
		args = append(args, repos[0])
	default:
		args = append(args, "--multiple")
		args = append(args, repos...)
	}
	cmd := exec.Command("git", args...)
	err := cmd.Run()
	switch er := err.(type) {
	case *exec.Error:
		return fmt.Errorf("%s %v", er.Name, er.Err)
	case *exec.ExitError:
		output, err := cmd.Output()
		return fmt.Errorf("%s %v", output, err)
	default:
		return err
	}
}

func (repo *Repo) Tags() ([]string, error) {
	if err := repo.FetchTags(); err != nil {
		return nil, err
	}
	cmd := exec.Command("git", "tag")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	go cmd.Run()
	var tags []string
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		tags = append(tags, scanner.Text())
	}
	return tags, nil
}

func (repo *Repo) HasOrigin() (bool, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return false, err
	}
	for _, remote := range remotes {
		if remote == "origin" {
			return true, nil
		}
	}
	return false, nil
}

func (repo *Repo) Remotes() ([]string, error) {
	args := []string{"remote"}
	cmd := exec.Command("git", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	remotes := []string{}
	go cmd.Run()
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return nil, err
		}
		remotes = append(remotes, scanner.Text())
	}
	return remotes, nil
}

func (repo *Repo) AddTag(tag semver.Version, comments ...string) error {
	if len(comments) == 0 {
		comments = []string{tag.String()}
	}
	message := strconv.QuoteToASCII(strings.Join(comments, "; "))
	cmd := exec.Command("git", "tag", "-a", "v"+tag.String(), "-m", message)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	switch err := cmd.Run().(type) {
	case *exec.ExitError:
		return fmt.Errorf("%s", buf.String())
	default:
		return err
	}
}
