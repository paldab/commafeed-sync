package commafeed

import (
	"fmt"
	"net/http"
)

func (api CFApi) SubscribeFeed(url, title, categoryId string) (*http.Response, error) {
	endpoint := "/rest/feed/subscribe"

	payload := RequestBody{
		"url":        url,
		"title":      title,
		"categoryId": categoryId,
	}

	resp, err := doPost(&api, endpoint, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to feed: %w", err)
	}

	defer resp.Body.Close()

	return resp, err
}

func (api CFApi) UnsubscribeFeed(id int) (*http.Response, error) {
	endpoint := "/rest/feed/unsubscribe"

	payload := RequestBody{
		"id": id,
	}

	resp, err := doPost(&api, endpoint, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to unsubscribe from feed: %w", err)
	}

	defer resp.Body.Close()

	return resp, err
}
