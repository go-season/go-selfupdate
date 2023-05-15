package selfupdate

import (
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

const (
	DefaultPrivateToken = "WTvPcYiim1W4322vpm9_"
	DefaultBaseURL      = "http://gitlab.idc.xiaozhu.com/api/v4"
)

type Updater struct {
	api *gitlab.Client
}

func DefaultUpdater() *Updater {
	token := os.Getenv("PRIVATE_TOKEN")
	if token == "" {
		token = DefaultPrivateToken
	}
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	apiClient := gitlab.NewClient(nil, token)
	err := apiClient.SetBaseURL(baseURL)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &Updater{
		api: apiClient,
	}
}
