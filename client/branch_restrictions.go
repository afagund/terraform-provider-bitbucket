package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetBranchRestrictions(repositorySlug string) (*Paginated[BranchRestriction], error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s/%s/branch-restrictions", c.Host, c.Workspace, repositorySlug), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var paginated Paginated[BranchRestriction]
	err = json.Unmarshal(body, &paginated)
	if err != nil {
		return nil, err
	}

	return &paginated, nil
}

func (c *Client) GetBranchRestriction(repositorySlug string, id int) (*BranchRestriction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s/%s/branch-restrictions/%d", c.Host, c.Workspace, repositorySlug, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var branchRestriction BranchRestriction
	err = json.Unmarshal(body, &branchRestriction)
	if err != nil {
		return nil, err
	}

	return &branchRestriction, nil
}

func (c *Client) CreateBranchRestriction(repositorySlug string, newBranchRestriction BranchRestriction) (*BranchRestriction, error) {
	rb, err := json.Marshal(newBranchRestriction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repositories/%s/%s/branch-restrictions", c.Host, c.Workspace, repositorySlug), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var branchRestriction BranchRestriction
	err = json.Unmarshal(body, &branchRestriction)
	if err != nil {
		return nil, err
	}

	return &branchRestriction, nil
}

func (c *Client) UpdateBranchRestriction(repositorySlug string, id int, newBranchRestriction BranchRestriction) (*BranchRestriction, error) {
	rb, err := json.Marshal(newBranchRestriction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/repositories/%s/%s/branch-restrictions/%d", c.Host, c.Workspace, repositorySlug, id), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var branchRestriction BranchRestriction
	err = json.Unmarshal(body, &branchRestriction)
	if err != nil {
		return nil, err
	}

	return &branchRestriction, nil
}

func (c *Client) DeleteBranchRestriction(repositorySlug string, id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/repositories/%s/%s/branch-restrictions/%d", c.Host, c.Workspace, repositorySlug, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
