package utils

import (
	"os"
	"paldab/commafeed-feed-sync/internal/models"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	CommafeedPathPrefix = "All/"
	CommafeedFeedSuffix = ".atom"
)

func LoadFeedsConfig(path string) (models.Config, error) {
	f, err := os.ReadFile(path)

	var feeds models.Config
	if err != nil {
		return feeds, err
	}

	err = yaml.Unmarshal(f, &feeds)

	return feeds, err
}

// Adding .atom to feed url because commafeed format
func CastToCommafeedFeedUrl(url string) string {
	var lookupCommafeedFeedUrl string = url
	isValidCommafeedUrl := strings.HasSuffix(lookupCommafeedFeedUrl, CommafeedFeedSuffix)

	if !isValidCommafeedUrl {
		lookupCommafeedFeedUrl = lookupCommafeedFeedUrl + CommafeedFeedSuffix
	}

	return lookupCommafeedFeedUrl
}
