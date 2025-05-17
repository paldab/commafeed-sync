package commafeed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"paldab/commafeed-feed-sync/internal/models"
	"time"
)

const (
	contentType           = "application/json"
	requestTimeoutSeconds = 10 * time.Second
)

type CommaFeedClient interface {
	GetCategories() (models.CommafeedCategoryResponse, error)
	CreateCategory(name, parentId string) (string, error)
	DeleteCategory(id int) (*http.Response, error)
	SubscribeFeed(url, title, categoryId string) (*http.Response, error)
	UnsubscribeFeed(id int) (*http.Response, error)
}

type RequestBody = map[string]interface{}

type CFApi struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Client   *http.Client
}

func (api CFApi) getUrl(endpoint string) string {
	// if !strings.HasPrefix(api.Url, "http") {
	// 	api.Url = fmt.Sprintf("http://%s", api.Url)
	// }

	return fmt.Sprintf("https://%s%s", api.Url, endpoint)
}

func doPost(api *CFApi, endpoint string, data RequestBody) (*http.Response, error) {
	reqUrl := fmt.Sprintf("%s%s", api.Url, endpoint)
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(jsonData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(api.Username, api.Password)

	// fmt.Println(data, "req")
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status: %d - %s", resp.StatusCode, respBody)
	}

	return resp, nil
}

func NewCFApi(url, username, password string) (*CFApi, error) {
	api := CFApi{
		Url:      url,
		Username: username,
		Password: password,
		Client:   &http.Client{},
	}

	// isValidCredentials := api.login()

	// if !isValidCredentials {
	// 	return &CFApi{}, fmt.Errorf("Supplied authentication credentials invalid")
	// }

	return &api, nil
}
