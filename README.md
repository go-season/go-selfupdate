# Self-Update Mechanism for Go Commands Using GitLab

 [go-gitlab-selfupdate](http://github.com/go-season/go-selfupdate) is a Go library that reference `go-github-selfupdate` to provide a self-update mechanism to command line tools.

## Install

```bash
go get -insecure github.com/go-season/go-selfupdate
```

## Usage

### Code Usage

It provides `selfupdate` package.

- `selfupdate.UpdateSelf()`: Detect the latest version of itself and run self update.
- `selfupdate.DetectLatest()`: Detect the latest version of given repository.
- `selfupdate.DetectVersion()`: Detect the user defined version of given repository.

Following is the easiest way to use this package.

```go
import (
    "log"
    "github.com/blang/semver"
    "github.com/go-season/go-selfupdate"
)

const version = "1.0.0"

func doSelfUpdate() {
    v := semver.MustParse(version)
    latest, err := selfupdate.UpdateSelf(v, "repoGroup/repoName")
    if err != nil {
        log.Println("Binary update failed:", err)
        return
    }
    if latest.Version.Equals(v) {
        log.Println("Current binary is the latest version", version)
    } else {
        log.Println("Successfully updated to version", latest.Version)
        log.Println("Release note:\n", latest.Description)
    }
}
```

## Dependencies

This library utilizes
- [go-gitlab](github.com/xanzy/go-gitlab) to retrieve the information of releases
- [go-update](github.com/inconshreveable/go-update) to replace current binary
- [semver](github.com/blang/semver) to compare versions