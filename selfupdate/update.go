package selfupdate

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
)

func (up *Updater) UpdateSelf(current semver.Version, slug string) (*Release, error) {
	cmdPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	return up.UpdateCommand(cmdPath, current, slug)
}

func (up *Updater) UpdateCommand(cmdPath string, current semver.Version, slug string) (*Release, error) {
	stat, err := os.Lstat(cmdPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to stat '%s'. File may not exist: %s", cmdPath, err)
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		p, err := filepath.EvalSymlinks(cmdPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to resolve symlink '%s' for executable: %s", cmdPath, err)
		}
		cmdPath = p
	}

	rel, ok, err := up.DetectLatest(slug)
	if err != nil {
		return nil, err
	}
	if !ok {
		log.Println("No release detected. Current version is considered up-to-date")
		return &Release{Version: current}, nil
	}
	if current.Equals(rel.Version) {
		log.Println("Current version", current, "is the latest. Update is not needed")
		return rel, nil
	}
	log.Println("Will update", cmdPath, "to the latest version", rel.Version)
	if err := up.UpdateTo(rel, cmdPath); err != nil {
		return nil, err
	}
	return rel, nil
}

func (up *Updater) UpdateTo(rel *Release, cmdPath string) error {
	src, err := up.DownloadReleaseAsset(rel.AssertURL)
	if err != nil {
		return err
	}

	return Update(src, rel.AssertURL, cmdPath)
}

func Update(src io.ReadCloser, assetURL, cmdPath string) error {
	defer src.Close()

	log.Println("Will update", cmdPath, "to the latest downloaded from", assetURL)
	return update.Apply(src, update.Options{
		TargetPath: cmdPath,
	})
}

func (up *Updater) DownloadReleaseAsset(assetURL string) (rc io.ReadCloser, err error) {
	resp, err := http.Get(assetURL)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func UpdateSelf(current semver.Version, slug string) (*Release, error) {
	return DefaultUpdater().UpdateSelf(current, slug)
}
