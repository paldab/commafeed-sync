package main

import (
	"log"
	"paldab/commafeed-feed-sync/config"
	"paldab/commafeed-feed-sync/internal/sync"
	"paldab/commafeed-feed-sync/internal/tree"
)

func main() {
	feedsConfig, api := config.Setup()
	commafeedDataTree, err := api.GetCategories()

	if err != nil {
		log.Fatal(err)
	}

	dataTree, err := tree.BuildTreeFromFeedsConfig(feedsConfig)

	if err != nil {
		log.Fatal(err)
	}

	sync.Sync(dataTree, commafeedDataTree, api)
}
