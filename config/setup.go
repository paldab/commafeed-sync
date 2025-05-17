package config

import (
	"log"
	"paldab/commafeed-feed-sync/internal/commafeed"
	"paldab/commafeed-feed-sync/internal/models"
	"paldab/commafeed-feed-sync/utils"
)

func Setup() (models.Config, *commafeed.CFApi) {
	url, err := GetCFUrl()

	if err != nil {
		log.Fatalf("Failed the setup %v", err)
	}

	configPath := GetFeedsConfigPath()
	feeds, err := utils.LoadFeedsConfig(configPath)

	if err != nil {
		log.Fatalf("Failed the setup %v", err)
	}

	username, password := GetCredentials()

	api, err := commafeed.NewCFApi(url, username, password)

	if err != nil {
		log.Fatalf("Failed the setup %v", err)
	}

	return feeds, api
}
