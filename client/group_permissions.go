package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetGroupPermissions(repositorySlug string) (*Paginated[GroupPermission], error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s/%s/permissions-config/groups", c.Host, c.Workspace, repositorySlug), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var paginated Paginated[GroupPermission]
	err = json.Unmarshal(body, &paginated)
	if err != nil {
		return nil, err
	}

	return &paginated, nil
}

func (c *Client) GetGroupPermission(repositorySlug, groupSlug string) (*GroupPermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s/%s/permissions-config/groups/%s", c.Host, c.Workspace, repositorySlug, groupSlug), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var groupPermission GroupPermission
	err = json.Unmarshal(body, &groupPermission)
	if err != nil {
		return nil, err
	}

	return &groupPermission, nil
}

func (c *Client) CreateGroupPermission(repositorySlug, groupSlug string, newGroupPermission GroupPermission) (*GroupPermission, error) {
	rb, err := json.Marshal(newGroupPermission)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/repositories/%s/%s/permissions-config/groups/%s", c.Host, c.Workspace, repositorySlug, groupSlug), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var groupPermission GroupPermission
	err = json.Unmarshal(body, &groupPermission)
	if err != nil {
		return nil, err
	}

	return &groupPermission, nil
}

func (c *Client) UpdateGroupPermission(repositorySlug, groupSlug string, newGroupPermission GroupPermission) (*GroupPermission, error) {
	rb, err := json.Marshal(newGroupPermission)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/repositories/%s/%s/permissions-config/groups/%s", c.Host, c.Workspace, repositorySlug, groupSlug), strings.NewReader(string(rb)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var groupPermission GroupPermission
	err = json.Unmarshal(body, &groupPermission)
	if err != nil {
		return nil, err
	}

	return &groupPermission, nil
}

func (c *Client) DeleteGroupPermission(repositorySlug, groupSlug string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/repositories/%s/%s/permissions-config/groups/%s", c.Host, c.Workspace, repositorySlug, groupSlug), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
