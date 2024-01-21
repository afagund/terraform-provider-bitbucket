package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetRepositories() (*Paginated[Repository], error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s", c.Host, c.Workspace), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var paginated Paginated[Repository]
	err = json.Unmarshal(body, &paginated)
	if err != nil {
		return nil, err
	}

	return &paginated, nil
}

func (c *Client) GetRepository(slug string) (*Repository, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s/%s", c.Host, c.Workspace, slug), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var repository Repository
	err = json.Unmarshal(body, &repository)
	if err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) CreateRepository(slug string, newRepository Repository) (*Repository, error) {
	rb, err := json.Marshal(newRepository)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repositories/%s/%s", c.Host, c.Workspace, slug), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var repository Repository
	err = json.Unmarshal(body, &repository)
	if err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) UpdateRepository(slug string, newRepository Repository) (*Repository, error) {
	rb, err := json.Marshal(newRepository)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/repositories/%s/%s", c.Host, c.Workspace, slug), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var repository Repository
	err = json.Unmarshal(body, &repository)
	if err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *Client) DeleteRepository(slug string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/repositories/%s/%s", c.Host, c.Workspace, slug), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
