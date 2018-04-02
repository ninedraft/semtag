package main

import (
	"errors"

	"github.com/blang/semver"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	ErrNoSemverTags = errors.New("semver tags not found")
)

func getTags(repPath string) ([]semver.Version, error) {
	rep, err := git.PlainOpen(repPath)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = rep.Fetch(&git.FetchOptions{RemoteName: "origin"})
	switch err {
	case nil, git.NoErrAlreadyUpToDate, git.ErrRemoteNotFound:
		// pass
	default:
		return nil, err
	}
	tags, err := rep.TagObjects()
	if err != nil {
		return nil, err
	}
	var stringTags []string
	err = tags.ForEach(func(tag *object.Tag) error {
		stringTags = append(stringTags, tag.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}
	localTag, err := readSemverFile()
	if err == nil {
		stringTags = append(stringTags, localTag.String())
	}
	semverTags := make([]semver.Version, 0, len(stringTags))
	for _, tag := range stringTags {
		vers, err := semver.ParseTolerant(tag)
		if err == nil {
			semverTags = append(semverTags, vers)
		}
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
