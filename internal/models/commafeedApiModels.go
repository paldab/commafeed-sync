package models

type CommafeedCategoryResponse struct {
	ID         string                       `json:"id"`
	ParentId   string                       `json:"parentId"`
	ParentName string                       `json:"parentName"`
	Name       string                       `json:"name"`
	Feeds      []CommafeedFeedResponse      `json:"feeds"`
	Expanded   bool                         `json:"expanded"`
	Position   int                          `json:"position"`
	Children   *[]CommafeedCategoryResponse `json:"children"`
}

type CommafeedFeedResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Message        string `json:"message"`
	ErrorCount     int    `json:"errorCount"`
	LastRefresh    int    `json:"lastRefresh"`
	NextRefresh    int    `json:"nextRefresh"`
	FeedUrl        string `json:"feedUrl"`
	FeedLink       string `json:"feedLink"`
	IconUrl        string `json:"iconUrl"`
	Unread         int    `json:"unread"`
	CategoryId     string `json:"categoryId"`
	Position       int    `json:"position"`
	NewestItemTime int    `json:"newestItemTime"`
	Filter         any    `json:"filter"`
}
