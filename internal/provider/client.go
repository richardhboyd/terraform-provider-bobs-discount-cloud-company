package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "https://api.us-east-1.whybobs.com"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient -
func NewClient(host, api_key *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		// Default Hashicups URL
		HostURL: HostURL,
		Token:   *api_key,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if api_key == nil {
		return &c, nil
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Set("api_key", token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) CreateDatabaseItem(createDatabaseItemRequest CreateDatabaseItemRequest, database_id string) (*CreateDatabaseItemResponse, error) {
	rb, err := json.Marshal(createDatabaseItemRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/database/%s/items", c.HostURL, database_id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	order := CreateDatabaseItemResponse{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// CreateDatabase - Create new order
func (c *Client) CreateDatabase(createDatabaseRequest CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	rb, err := json.Marshal(createDatabaseRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/database", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	order := CreateDatabaseResponse{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetDatabase - Get database
func (c *Client) GetDatabase(database_id string) (*CreateDatabaseResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/database/%s", c.HostURL, database_id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	database := CreateDatabaseResponse{}
	err = json.Unmarshal(body, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

// ListDatabases
func (c *Client) ListDatabases() (*ListDatabasesResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/database", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	order := ListDatabasesResponse{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// DeleteDatabase - Deletes a database
func (c *Client) DeleteDatabase(database_id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/database/%s", c.HostURL, database_id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// Database -
type Database struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateDatabaseRequest struct {
	Name string `json:"name"`
}

type GetDatabaseRequest struct {
	Id string `json:"id"`
}

type GetDatabaseResponse struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

// type ListDatabasesResponse []Database

type ListDatabasesResponse struct {
	Databases []Database `json:"databases"`
}
type CreateDatabaseResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// OrderItem -
type DatabaseItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateDatabaseItemRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateDatabaseItemResponse []DatabaseItem

// type CreateDatabaseItemResponse struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }
