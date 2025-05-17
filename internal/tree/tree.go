package tree

import (
	"paldab/commafeed-feed-sync/internal/models"
	"paldab/commafeed-feed-sync/utils"
	"strings"
)

// Node = category

// has to be an array because there are multiple categories that are not nested so a list of trees
func BuildTreeFromFeedsConfig(config models.Config) ([]*models.Category, error) {
	if len(config.CommafeedSetup) == 0 {
		return nil, ErrEmptyConfig
	}

	var root []*models.Category

	for _, c := range config.CommafeedSetup {
		nameParts := strings.Split(c.Name, "/")
		InsertIntoTree(nameParts, c.Feeds, &root, "")
	}

	return root, nil
}

func InsertIntoTree(names []string, feeds []models.Feed, tree *[]*models.Category, parentName string) {
	if len(names) == 0 {
		return
	}

	currentName := strings.TrimSpace(names[0])

	var child *models.Category

	// Does node exist
	for _, category := range *tree {
		if currentName == category.Name {
			child = category
			break
		}
	}

	// Node doesn not exist as child
	if child == nil {
		child = &models.Category{
			Name:       currentName,
			ParentName: &parentName,
			Children:   []*models.Category{},
		}
		*tree = append(*tree, child)
	}

	// Filter disabled feeds
	enabledFeeds := []models.Feed{}
	for _, feed := range feeds {
		if !feed.Disabled {
			// Cast already to commafeedFormat
			feed.Url = utils.CastToCommafeedFeedUrl(feed.Url)
			enabledFeeds = append(enabledFeeds, feed)
		}
	}

	// If last category
	if len(names) == 1 {
		child.Feeds = enabledFeeds
	} else {
		InsertIntoTree(names[1:], enabledFeeds, &child.Children, currentName)
	}
}
