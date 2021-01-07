package appledevconn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
)

// EnsureConnection ...
func EnsureConnection() error {
	buildURL, buildAPIToken := os.Getenv("BITRISE_BUILD_URL"), os.Getenv("BITRISE_BUILD_API_TOKEN")
	if buildURL == "" || buildAPIToken == "" {
		log.Warnf("BITRISE_BUILD_URL and/or BITRISE_BUILD_API_TOKEN envs are not set")
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/apple_developer_portal_data.json", buildURL), nil)
	if err != nil {
		return err
	}

	req.Header.Add("BUILD_API_TOKEN", buildAPIToken)

	body, err := performRequest(req)
	if err != nil {
		return err
	}

	fmt.Println(body)

	return nil
}

func performRequest(req *http.Request) ([]byte, error) {
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		// On error, any Response can be ignored
		return nil, fmt.Errorf("failed to perform request, error: %s", err)
	}

	// The client must close the response body when finished with it
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Warnf("Failed to close response body, error: %s", cerr)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, error: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return body, fmt.Errorf("GET %s failed with status code: %d", req.URL.String(), response.StatusCode)
	}

	return body, nil
}
