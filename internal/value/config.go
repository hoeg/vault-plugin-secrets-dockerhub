package value

import (
	"fmt"
	"slices"
	"time"
)

const DefaultTTL time.Duration = 5 * time.Minute

// Config holds to values needed to issue a new Docker Hub access token
type Config struct {
	Scopes   []string      `json:"scopes"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
}

var validScopes = []string{"admin", "write", "read", "public_read"}

func NewConfig(username, password string, scopes []string) (*Config, error) {
	for _, s := range scopes {
		if !slices.Contains(validScopes, s) {
			return nil, fmt.Errorf("invalid slice provided: %s", s)
		}
	}
	return &Config{
		Scopes:   scopes,
		Username: username,
		Password: password,
		TTL:      DefaultTTL,
	}, nil
}
