package selfupdate

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/xanzy/go-gitlab"
)

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

func (up *Updater) DetectLatest(slug string) (release *Release, found bool, err error) {
	return up.DetectVersion(slug, "")
}

func (up *Updater) DetectVersion(slug string, version string) (release *Release, found bool, err error) {
	repo := strings.Split(slug, "/")

	rels, res, err := up.api.Releases.ListReleases(slug, nil, nil)
	if err != nil {
		log.Println("API returned an error response:", err)
		if res != nil && res.StatusCode == 404 {
			// 404 means repository not found or release not found. It's not an error here.
			err = nil
			log.Println("API returned 404. Repository or release not found")
		}
		return nil, false, err
	}

	project, _, err := up.api.Projects.GetProject(slug, nil, nil)
	if err != nil {
		log.Println("Get project failed.")
		return nil, false, err
	}

	rel, asset, ver, found := findReleaseAndAsset(repo[1], rels, version)
	if !found {
		return nil, false, nil
	}

	url := asset.URL
	//log.Println("Successfully fetched the latest release. tag:", rel.TagName, ", name:", rel.Name, ", URL:", project.WebURL, ", Asset:", url)

	release = &Release{
		Version:     ver,
		AssertURL:   url,
		URL:         project.WebURL,
		Name:        rel.Name,
		RepoOwner:   repo[0],
		RepoName:    repo[1],
		Description: rel.Description,
	}

	return release, true, nil
}

func findReleaseAndAsset(repoName string, rels []*gitlab.Release, targetVersion string) (*gitlab.Release, *gitlab.ReleaseLink, semver.Version, bool) {
	suffixes := make([]string, 0, 2*7*2)
	for _, sep := range []rune{'_', '-'} {
		suffix := fmt.Sprintf("%s%c%s%c%s", repoName, sep, runtime.GOOS, sep, runtime.GOARCH)
		suffixes = append(suffixes, suffix)
	}

	var ver semver.Version
	var asset *gitlab.ReleaseLink
	var release *gitlab.Release

	for _, rel := range rels {
		if a, v, ok := findAssetFromRelease(rel, suffixes, targetVersion); ok {
			if release == nil || v.GTE(ver) {
				ver = v
				asset = a
				release = rel
			}
		}
	}
	if release == nil {
		log.Println("Could not find any release for", runtime.GOOS, "and", runtime.GOARCH)
		return nil, nil, semver.Version{}, false
	}

	return release, asset, ver, true
}

func findAssetFromRelease(rel *gitlab.Release, suffixes []string, targetVersion string) (*gitlab.ReleaseLink, semver.Version, bool) {
	if targetVersion != "" && targetVersion != rel.TagName {
		//log.Println("Skip", rel.TagName, "not matching to specified version", targetVersion)
		return nil, semver.Version{}, false
	}

	verText := rel.TagName
	indices := reVersion.FindStringIndex(verText)
	if indices == nil {
		//log.Println("Skip version not adopting semver", verText)
		return nil, semver.Version{}, false
	}
	if indices[0] > 0 {
		//log.Println("Strip prefix of version", verText[:indices[0]], "from", verText)
		verText = verText[indices[0]:]
	}

	ver, err := semver.Make(verText)
	if err != nil {
		//log.Println("Failed to parse a semantic version", verText)
		return nil, semver.Version{}, false
	}

	for _, asset := range rel.Assets.Links {
		name := asset.URL
		for _, s := range suffixes {
			if strings.HasSuffix(name, s) {
				return asset, ver, true
			}
		}
	}

	//log.Println("No suitable asset was found in release", rel.TagName)
	return nil, semver.Version{}, false
}

func DetectLatest(slug string) (*Release, bool, error) {
	return DefaultUpdater().DetectLatest(slug)
}
