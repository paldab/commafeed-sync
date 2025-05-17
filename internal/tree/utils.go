package tree

import (
	"fmt"
	"log"
	"paldab/commafeed-feed-sync/internal/models"
	"strconv"
	"strings"
)

func getParentWithNode(node, parent string) string {
	const separator = "/"
	return strings.Trim(parent+separator+node, separator)
}

func FlattenDeclaredTree(tree []*models.Category, parent string, out map[string]*models.Category) {
	for _, node := range tree {
		currentPath := getParentWithNode(node.Name, parent)
		out[currentPath] = node
		if node.Children != nil {
			FlattenDeclaredTree(node.Children, currentPath, out)
		}
	}
}

func FlattenCommafeedTree(tree models.CommafeedCategoryResponse, parent string, out map[string]models.CommafeedCategoryResponse) {
	currentPath := getParentWithNode(tree.Name, parent)
	_, err := strconv.Atoi(tree.ID)

	if tree.ID != "all" && err != nil {
		log.Printf("Skipping invalid category ID %q at path %q", tree.ID, currentPath)
		return
	}

	out[currentPath] = tree

	if tree.Children != nil {
		for _, child := range *tree.Children {
			FlattenCommafeedTree(child, currentPath, out)
		}
	}
}

// For Dev
func PrintTree(categories []*models.Category, indent int) {
	for _, cat := range categories {
		fmt.Printf("%s- %s (parent: %s)\n", strings.Repeat("  ", indent), cat.Name, *cat.ParentName)
		for _, feed := range cat.Feeds {
			fmt.Printf("%s  📄 %s (%s)\n", strings.Repeat("  ", indent+1), feed.Name, feed.Url)
		}
		PrintTree(cat.Children, indent+1)
	}
}

// type CategoryWithDepth struct {
// 	Category models.Category
// 	Depth    int
// }
//
// // Using DFS UNUSED
// func GetDeepestCategory(category models.Category) (models.Category, int) {
// 	var deepestPoint models.Category = models.Category{}
// 	maxDepth := 1 // Root is 1
//
// 	stack := []CategoryWithDepth{
// 		{Category: category, Depth: maxDepth},
// 	}
//
// 	// visited := []models.Category{category}
//
// 	for len(stack) > 0 {
// 		lastElement := len(stack) - 1
// 		curr := stack[lastElement]
//
// 		// Remove curr from stack
// 		stack = stack[:lastElement]
//
// 		// Compare if current depth is deeper than maxDepth
// 		if curr.Depth > maxDepth {
// 			maxDepth = curr.Depth
// 			deepestPoint = curr.Category
// 		}
//
// 		for _, child := range curr.Category.Children {
// 			stack = append(stack, CategoryWithDepth{Category: *child, Depth: curr.Depth + 1})
// 		}
//
// 	}
//
// 	return deepestPoint, maxDepth
// }
