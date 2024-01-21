package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const Host string = "https://api.bitbucket.org/2.0"

type Client struct {
	Host       string
	Workspace  string
	Token      string
	HTTPClient *http.Client
}

type Error struct {
	Message string `json:"message"`
}

type ResponseErr struct {
	ErrorDetails Error `json:"error"`
}

type StatusErr struct {
	StatusCode int
	Message    string
}

func (s StatusErr) Error() string {
	return s.Message
}

func NewClient(host, workspace, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Host:       Host,
	}

	if host != nil {
		c.Host = *host
	}

	if workspace != nil {
		c.Workspace = *workspace
	}

	if token != nil {
		c.Token = *token
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Authorization", "Basic "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		var responseErr ResponseErr
		err = json.Unmarshal(body, &responseErr)
		if err != nil {
			return nil, fmt.Errorf("the Bitbucket API returned, status: %d", res.StatusCode)
		}

		if res.StatusCode == 404 {
			return nil, StatusErr{
				StatusCode: res.StatusCode,
				Message:    fmt.Sprintf("the Bitbucket API returned, status: %d message: %s", res.StatusCode, responseErr.ErrorDetails.Message),
			}
		}

		return nil, fmt.Errorf("the Bitbucket API returned, status: %d message: %s", res.StatusCode, responseErr.ErrorDetails.Message)
	}

	return body, err
}
