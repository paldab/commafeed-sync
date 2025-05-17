package tree

import (
	"fmt"
	"paldab/commafeed-feed-sync/internal/models"
	"paldab/commafeed-feed-sync/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConfig(t *testing.T) {
	data, err := utils.LoadFeedsConfig("../../testdata/empty.yaml")
	assert.Nil(t, err)

	categories, err := BuildTreeFromFeedsConfig(data)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyConfig, err)
	assert.Len(t, categories, 0)
}

func TestEmptyConfigWithSetupKey(t *testing.T) {
	data, err := utils.LoadFeedsConfig("../../testdata/empty_config.yaml")
	assert.Nil(t, err)

	categories, err := BuildTreeFromFeedsConfig(data)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyConfig, err)
	assert.Len(t, categories, 0)
}

func TestBuildTreeRegular(t *testing.T) {
	data, err := utils.LoadFeedsConfig("../../testdata/feed_valid.yaml")
	assert.Nil(t, err)

	categories, err := BuildTreeFromFeedsConfig(data)

	assert.Nil(t, err)
	assert.IsType(t, []*models.Category{}, categories)
	assert.Len(t, categories, 1)

	category := categories[0]

	expectedCategoryName := "Unnested"
	assert.Equal(t, expectedCategoryName, category.Name)
	assert.Len(t, category.Feeds, 1)

	feed := category.Feeds[0]

	assert.Equal(t, feed.Name, "Commafeed Releases")
	assert.Len(t, category.Children, 0)
}

func TestBuildTreeMultiple(t *testing.T) {
	data, err := utils.LoadFeedsConfig("../../testdata/multi_feeds_valid.yaml")
	assert.Nil(t, err)

	categories, err := BuildTreeFromFeedsConfig(data)

	assert.Nil(t, err)
	assert.IsType(t, []*models.Category{}, categories)
	assert.Len(t, categories, 4)

	// excluding disabled
	expectedFeedCounts := []int{1, 1, 2, 2}

	for i, category := range categories {
		expectedCategoryName := fmt.Sprintf("Test%d", i+1)

		assert.Equal(t, expectedCategoryName, category.Name)
		assert.Len(t, category.Children, 0)
		assert.Len(t, category.Feeds, expectedFeedCounts[i])
	}
}

// func TestBuildTreeNested(t *testing.T) {
// 	data, err := utils.LoadFeedsConfig("../../testdata/feeds_valid_advanced")
// 	assert.Nil(t, err)
//
// 	categories, err := BuildTreeFromFeedsConfig(data)
// 	assert.Nil(t, err)
// 	assert.Len(t, categories, 5)
//
// 	expected := []struct {
// 		Feeds     int
// 		Depth     int
// 		DepthName string
// 	}{
// 		{
// 			Feeds:     1,
// 			Depth:     1,
// 			DepthName: "Releases",
// 		},
// 		{
// 			Feeds:     2,
// 			Depth:     1,
// 			DepthName: "Kubernetes",
// 		},
// 		{
// 			Feeds:     3,
// 			Depth:     2,
// 			DepthName: "Software Releases",
// 		},
// 		{
// 			Feeds:     4,
// 			Depth:     2,
// 			DepthName: "Software Releases",
// 		},
// 		{
// 			Feeds:     4,
// 			Depth:     3,
// 			DepthName: "Interests",
// 		},
// 	}
//
// 	for i, category := range categories {
// 		deepestCategory, depth := utils.GetDeepestCategory(category)
//
// 		assert.Equal(t, deepestCategory.Name, expected[i].DepthName)
// 		assert.Equal(t, deepestCategory.Feeds, expected[i].Feeds)
// 		assert.Equal(t, depth, expected[i].Depth)
// 	}
// }
