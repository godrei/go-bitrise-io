package appledevconn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/xcode-project/pretty"
)

// Cookie ...
type Cookie struct {
	Name       string     `json:"name,omitempty"`
	Value      string     `json:"value,omitempty"`
	Domain     string     `json:"domain,omitempty"`
	ForDomain  bool       `json:"for_domain,omitempty"`
	Path       string     `json:"path,omitempty"`
	Secure     bool       `json:"secure,omitempty"`
	HTTPOnly   bool       `json:"httponly,omitempty"`
	Expires    *time.Time `json:"expires,omitempty"`
	MaxAge     int        `json:"max_age,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	AccessedAt *time.Time `json:"accessed_at,omitempty"`
}

// SessionConnection ...
type SessionConnection struct {
	AppleID              string              `json:"apple_id"`
	Password             string              `json:"password"`
	ConnectionExpiryDate string              `json:"connection_expiry_date"`
	SessionCookie        map[string][]Cookie `json:"session_cookies"`
}

// JWTConnection ...
type JWTConnection struct {
	KeyID      string `json:"key_id"`
	IssuerID   string `json:"issuer_id"`
	PrivateKey string `json:"private_key"`
}

// Device ...
type Device struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	DeviceID   string `json:"device_identifier"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeviceType string `json:"device_type"`
}

// AppleDeveloperConnection ...
type AppleDeveloperConnection struct {
	SessionConnection
	JWTConnection
	Devices []Device `json:"test_devices"`
}

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

	var conn AppleDeveloperConnection
	if err := json.Unmarshal(body, &conn); err != nil {
		return err
	}

	sessionConn := conn.SessionConnection
	fmt.Printf("sessionConn:\n%s\n", pretty.Object(sessionConn))

	jwtConn := conn.JWTConnection
	fmt.Printf("jwtConn:\n%s\n", pretty.Object(jwtConn))

	devices := conn.Devices
	fmt.Printf("devices:\n%s\n", pretty.Object(devices))

	// fmt.Println(string(body))

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
