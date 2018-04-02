package main

import (
	"errors"

	"github.com/blang/semver"
	"github.com/ninedraft/semtag/git"
)

var (
	ErrNoSemverTags = errors.New("semver tags not found")
)

func getTags(repPath string) ([]semver.Version, error) {
	repo := git.Repo{}
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	semverTags := make([]semver.Version, 0, len(tags))
	for _, tag := range tags {
		if version, err := semver.ParseTolerant(tag); err == nil {
			semverTags = append(semverTags, version)
		}
	}
	localTag, err := getLocalVersion()
	if err == nil {
		semverTags = append(semverTags, localTag)
	}
	semver.Sort(semverTags)
	return semverTags, nil
}

func getHeadTag(repPath string) (semver.Version, error) {
	tags, err := getTags(repPath)
	if err != nil {
		return semver.Version{}, err
	}
	if len(tags) > 0 {
		return tags[0], nil
	}
	return semver.Version{}, ErrNoSemverTags
}
