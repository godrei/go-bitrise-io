package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client ...
type Client struct {
	BaseURL *url.URL

	client   *http.Client
	apiToken string
}

// NewClient ...
func NewClient(buildURL, buildAPIToken string) (*Client, error) {
	baseURL, err := url.Parse(buildURL)
	if err != nil {
		return nil, fmt.Errorf("invalid build url (%s): %s", buildURL, err)
	}

	return &Client{
		BaseURL:  baseURL,
		client:   http.DefaultClient,
		apiToken: buildAPIToken,
	}, nil
}

func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing endpoint failed: %s", err)
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %s", err)
	}

	req.Header.Add("BUILD_API_TOKEN", c.apiToken)

	return req, nil
}

func checkResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	msg := fmt.Sprintf("%s %s status code: %d", r.Request.Method, r.Request.URL.String(), r.StatusCode)

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		msg += fmt.Sprintf(", response: %s", string(data))
	}

	return errors.New(msg)
}

// Do ...
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("perform request failed: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Failed to close response body: %s\n", cerr)
		}
	}()

	if err := checkResponse(resp); err != nil {
		return err
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}

	return err
}
