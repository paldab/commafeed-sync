package sync

import (
	"fmt"
	"log"
	"paldab/commafeed-feed-sync/internal/commafeed"
	"paldab/commafeed-feed-sync/internal/models"
	"paldab/commafeed-feed-sync/internal/tree"
	"paldab/commafeed-feed-sync/utils"
	"strconv"
	"strings"
)

func Sync(dataTree []*models.Category, commafeedDataTree models.CommafeedCategoryResponse, client commafeed.CommaFeedClient) {
	declaredMap := map[string]*models.Category{}
	commafeedMap := map[string]models.CommafeedCategoryResponse{}

	tree.FlattenDeclaredTree(dataTree, "", declaredMap)
	tree.FlattenCommafeedTree(commafeedDataTree, "", commafeedMap)

	// Create or update feeds
	ascDeclaredPaths := sortMapKeysAscendingHierachy(declaredMap)

	for _, path := range ascDeclaredPaths {
		declaredCategory := declaredMap[path]
		finalID, err := syncCategories(path, declaredMap, commafeedMap, client)
		if err != nil {
			log.Printf("Failed to ensure category path %s: %v\n", path, err)
			continue
		}

		commafeedCategory, exists := lookupCommafeedMap(path, commafeedMap)
		if exists {
			syncFeeds(declaredCategory.Feeds, commafeedCategory.Feeds, finalID, client)
		}
	}

	// Sort map to delete from bottom to top to delete in right order
	descCommafeedPaths := sortMapKeysDescendingHierachy(commafeedMap)

	// Delete categories and feeds
	for _, path := range descCommafeedPaths {
		cat := commafeedMap[path]
		// Skip root
		if cat.ID == "all" {
			continue
		}

		if cat.ID == "0" {
			log.Printf("something went wrong with getting ID from commafeed. path: %s", path)
			continue
		}

		// if data is declared, don't delete
		if _, exists := lookupDeclaredMap(path, declaredMap); exists {
			continue
		}

		// Delete feeds
		for _, feed := range cat.Feeds {
			log.Printf("Removing feed with url: '%s' because not declared\n", feed.FeedUrl)
			_, err := client.UnsubscribeFeed(feed.ID)

			if err != nil {
				log.Printf("Failed to unsubscribe from feed: %s when deleting category: %s\n", feed.Name, cat.Name)
			}
		}

		intCatID, err := strconv.Atoi(cat.ID)

		if err != nil {
			log.Printf("Failed to cast commafeed category ID: %s to int. err: %v\n", cat.ID, err)
		}

		if _, err = client.DeleteCategory(intCatID); err != nil {
			log.Printf("Failed to delete category with ID: %d and name: %s. err: %v\n", intCatID, cat.Name, err)
		}
	}
}

func syncCategories(path string, declaredMap map[string]*models.Category, commafeedMap map[string]models.CommafeedCategoryResponse, client commafeed.CommaFeedClient) (string, error) {
	var fullPath string
	var parentID string
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if i == 0 {
			fullPath = utils.CommafeedPathPrefix + part
		} else {
			fullPath = fullPath + "/" + part
		}

		if cat, exists := commafeedMap[fullPath]; exists {
			parentID = cat.ID
			continue
		}

		// Create missing category
		log.Printf("Creating category: %q with parentID: %q (raw: %v)", part, parentID, []byte(parentID))
		newID, err := client.CreateCategory(part, parentID)
		if err != nil {
			return "", fmt.Errorf("failed to create category %s: %w", fullPath, err)
		}
		log.Printf("Created category %q with returned ID: %q (raw: %v)", part, newID, []byte(newID))

		declaredCat, exists := lookupDeclaredMap(fullPath, declaredMap)
		if !exists || declaredCat == nil {
			log.Printf("WARNING: declaredMap[%q] is nil", fullPath)
			continue // or return an error
		}

		newFeeds := declaredCat.Feeds
		for _, feed := range newFeeds {
			log.Printf("subscribing to new feed: %s at url: %s\n", feed.Name, feed.Url)
			client.SubscribeFeed(feed.Url, feed.Name, newID)
		}

		// Updates commafeedMap with parentIDS
		commafeedMap[fullPath] = models.CommafeedCategoryResponse{
			ID:       newID,
			Name:     part,
			ParentId: parentID,
			Feeds:    castToCommafeedFeeds(newFeeds),
		}

		parentID = newID
	}

	return parentID, nil
}

// compare feeds, if feed should not be there, delete, otherwise add
func syncFeeds(declaredFeeds []models.Feed, commafeedFeeds []models.CommafeedFeedResponse, categoryID string, client commafeed.CommaFeedClient) {
	commafeedMap := make(map[string]models.CommafeedFeedResponse)
	for _, feed := range commafeedFeeds {
		commafeedMap[feed.FeedUrl] = feed
	}

	declaredMap := make(map[string]models.Feed)
	for _, feed := range declaredFeeds {
		declaredMap[feed.Url] = feed
	}

	// CHeck if already exist...
	for url, feed := range declaredMap {
		if _, exists := commafeedMap[url]; !exists {
			log.Printf("subscribing to new feed: %s at url: %s\n", feed.Name, url)
			_, err := client.SubscribeFeed(feed.Url, feed.Name, categoryID)

			if err != nil {
				log.Printf("Failed to subscribe to feed %s: %v\n", url, err)
			}
		}
	}

	for url, feed := range commafeedMap {
		if _, exists := declaredMap[url]; !exists {
			_, err := client.UnsubscribeFeed(feed.ID)

			if err != nil {
				log.Printf("Failed to unsubscribe to feed %s: %v\n", url, err)
			} else {
				log.Printf("Unsubscribed feed: %s at url: %s\n", feed.Name, url)
			}
		}
	}
}
