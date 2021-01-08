package connections

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const appleDevPortalDataPath = "apple_developer_portal_data.json"

// Cookie ...
// https://github.com/bitrise-io/apple-dev-portal-api/blob/master/models/serializablecookie.go#L9
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
	SessionCookies       map[string][]Cookie `json:"session_cookies"`
}

// FastlaneSession ...
func (c SessionConnection) FastlaneSession() (string, error) {
	var convertedCookies []string
	var errs []string

	for _, cookies := range c.SessionCookies {
		for _, cookie := range cookies {
			if convertedCookies == nil {
				convertedCookies = append(convertedCookies, "---"+"\n")
			}

			tmpl, err := template.New("").Parse(`- !ruby/object:HTTP::Cookie
	  name: {{.Name}}
	  value: {{.Value}}
	  domain: {{.Domain}}
	  for_domain: {{.ForDomain}}
	  path: "{{.Path}}"
	`)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Failed to create golang template for the cookie: %v", c))
				continue
			}

			var b bytes.Buffer
			err = tmpl.Execute(&b, cookie)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Failed to parse cookie: %v", c))
				continue
			}

			convertedCookies = append(convertedCookies, b.String()+"\n")
		}
	}

	if len(errs) > 0 {
		return "", errors.New(strings.Join(errs, "\n"))
	}
	return strings.Join(convertedCookies, ""), nil
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
	SessionConnection *SessionConnection
	JWTConnection     *JWTConnection
	Devices           []Device
}

// GetAppleDeveloperConnection ...
func (c *Client) GetAppleDeveloperConnection() (*AppleDeveloperConnection, error) {
	req, err := c.newRequest(http.MethodGet, appleDevPortalDataPath)
	if err != nil {
		return nil, err
	}

	type data struct {
		*SessionConnection
		*JWTConnection
		Devices []Device `json:"test_devices"`
	}

	var conn data
	if err := c.do(req, &conn); err != nil {
		return nil, err
	}

	return &AppleDeveloperConnection{
		SessionConnection: conn.SessionConnection,
		JWTConnection:     conn.JWTConnection,
		Devices:           conn.Devices,
	}, nil
}
