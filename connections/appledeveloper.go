package connections

import (
	"net/http"
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
