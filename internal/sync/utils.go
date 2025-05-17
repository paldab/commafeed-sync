package sync

import (
	"paldab/commafeed-feed-sync/internal/models"
	"paldab/commafeed-feed-sync/utils"
	"sort"
	"strings"
)

func castToCommafeedFeeds(feeds []models.Feed) []models.CommafeedFeedResponse {
	var targetFeeds []models.CommafeedFeedResponse

	for _, feed := range feeds {
		targetFeeds = append(targetFeeds, models.CommafeedFeedResponse{
			Name:    feed.Name,
			FeedUrl: feed.Url,
		})
	}

	return targetFeeds
}

// Commafeed need prefixed path with /All
func lookupDeclaredMap[V any](path string, declaredMap map[string]*V) (*V, bool) {
	lookupPath := strings.TrimPrefix(path, utils.CommafeedPathPrefix)
	category, exists := declaredMap[lookupPath]
	return category, exists
}

func lookupCommafeedMap[V any](path string, commafeedMap map[string]V) (V, bool) {
	lookupPath := path
	if !strings.HasPrefix(path, utils.CommafeedPathPrefix) {
		lookupPath = utils.CommafeedPathPrefix + path
	}

	category, exists := commafeedMap[lookupPath]
	return category, exists
}

func sortMapKeysDescendingHierachy[V any](mapToSort map[string]V) []string {
	keys := make([]string, 0, len(mapToSort))
	for k := range mapToSort {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.Count(keys[i], "/") > strings.Count(keys[j], "/")
	})

	return keys
}

func sortMapKeysAscendingHierachy[V any](mapToSort map[string]V) []string {
	keys := make([]string, 0, len(mapToSort))
	for k := range mapToSort {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.Count(keys[i], "/") < strings.Count(keys[j], "/")
	})

	return keys
}
