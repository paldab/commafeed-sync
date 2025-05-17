package commafeed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"paldab/commafeed-feed-sync/internal/models"
	"strconv"
)

func (api CFApi) GetCategories() (models.CommafeedCategoryResponse, error) {
	endpoint := "/rest/category/get"
	fullUrl := fmt.Sprintf("%s%s", api.Url, endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeoutSeconds)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)

	if err != nil {
		return models.CommafeedCategoryResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(api.Username, api.Password)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return models.CommafeedCategoryResponse{}, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CommafeedCategoryResponse{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.CommafeedCategoryResponse{}, err
	}

	var data models.CommafeedCategoryResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return models.CommafeedCategoryResponse{}, err
	}

	return data, nil
}

func (api CFApi) CreateCategory(name, parentId string) (string, error) {
	endpoint := "/rest/category/add"

	payload := RequestBody{
		"name": name,
	}

	if parentId != "" {
		payload["parentId"] = parentId
	}

	resp, err := doPost(&api, endpoint, payload)
	if err != nil {
		return "", fmt.Errorf("create category request failed: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body. err: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	var resultID int

	if err = json.Unmarshal(respBody, &resultID); err != nil {
		return "", fmt.Errorf("failed to decode response: %w; body: %s", err, string(respBody))
	}

	return strconv.Itoa(resultID), nil
}

func (api CFApi) DeleteCategory(id int) (*http.Response, error) {
	endpoint := "/rest/category/delete"

	payload := RequestBody{
		"id": id,
	}

	resp, err := doPost(&api, endpoint, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}
	defer resp.Body.Close()

	return resp, err
}
