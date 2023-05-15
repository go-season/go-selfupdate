package selfupdate

import (
	"github.com/blang/semver"
)

type Release struct {
	Version semver.Version

	AssertURL string

	URL string

	Name string

	RepoOwner string

	RepoName string

	Description string
}
