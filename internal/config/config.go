package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathPatternConfig is the string used to define the base path of the config
// endpoint
var pathPatternConfig = fmt.Sprintf("config/%s", framework.GenericNameRegex(configUsername))

const configListKey = "dockerhub/config"

const (
	fmtErrConfRetrieval = "failed to get configuration from storage"
	fmtErrConfMarshal   = "failed to marshal configuration to JSON"
	fmtErrConfUnmarshal = "failed to unmarshal configuration from JSON"
	fmtErrConfPersist   = "failed to persist configuration to storage"
	fmtErrConfDelete    = "failed to delete configuration from storage"
)

const (
	configNamespace     = "namespace"
	descConfigNamespace = "Docker Hub namespace that should be configured."
	configUsername      = "username"
	descConfigUsername  = "Docker Hub username that will issue access tokens."
	configPassword      = "password"
	descConfigPassword  = "Password for the Docker Hub user."
)

const pathConfigHelpSyn = `
Configure the Docker Hub secrets plugin.
`

const DefaultTTL time.Duration = 5 * time.Minute

var pathConfigHelpDesc = ""

func Paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				configNamespace: {
					Type:        framework.TypeCommaStringSlice,
					Description: descConfigNamespace,
				},
				configUsername: {
					Type:        framework.TypeString,
					Description: descConfigUsername,
				},
				configPassword: {
					Type:        framework.TypeString,
					Description: descConfigPassword,
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated keys. If <= 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time a service account key is valid for. If <= 0, will use system default.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: handleCreate,
					Summary:  "Create a configuration for Docker Hub.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: handleDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: handleRead,
					Summary:  "Read the configuration for Docker Hub access tokens for a specific user.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: handleCreate,
					Summary:  "Update an existing configuration for Docker Hub.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: handleList,
					Summary:  "List all configurations for the Docker Hub engine.",
				},
			},
			HelpSynopsis:    pathConfigHelpSyn,
			HelpDescription: pathConfigHelpDesc,
		},
	}
}

// Config holds to values needed to issue a new Docker Hub access token
type Config struct {
	Namespace []string      `json:"namespace"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	TTL       time.Duration `json:"ttl"`
	MaxTTL    time.Duration `json:"max_ttl"`
}

func RetrieveConfig(ctx context.Context, username string, s logical.Storage) (*Config, error) {
	c := Config{}
	entry, err := s.Get(ctx, getStorePath(username))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfRetrieval, err)
	}

	if entry == nil || len(entry.Value) == 0 {
		return nil, fmt.Errorf("unable to finde configuration for %q", username)
	}

	if err := entry.DecodeJSON(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfUnmarshal, err)
	}

	return &c, nil
}

func handleCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	c := Config{}
	c.Username = data.Get(configUsername).(string)
	c.Password = data.Get(configPassword).(string)
	c.Namespace = data.Get(configNamespace).([]string)
	c.TTL = DefaultTTL

	ce, err := logical.StorageEntryJSON(getStorePath(c.Username), c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfMarshal, err)
	}
	if err = req.Storage.Put(ctx, ce); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfPersist, err)
	}
	return nil, nil
}

func handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := data.Get(configUsername).(string)
	return nil, req.Storage.Delete(ctx, getStorePath(u))
}

func getStorePath(u string) string {
	return fmt.Sprintf("dockerhub/config/%s", u)
}

func handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := data.Get(configUsername).(string)
	c, err := RetrieveConfig(ctx, u, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]interface{})

	if v := c.Username; v != "" {
		resp["username"] = v
	}
	if v := c.Namespace; v != nil {
		resp["namespace"] = v
	}
	resp["ttl"] = c.TTL.String()

	return &logical.Response{
		Data: resp,
	}, nil
}

func handleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configs, err := req.Storage.List(ctx, getStorePath(""))
	if err != nil {
		return nil, err
	}

	resp := make(map[string]interface{})
	for _, sc := range configs {
		c := Config{}
		if err := json.Unmarshal([]byte(sc), &c); err != nil {
			return nil, fmt.Errorf("%s: %w", fmtErrConfUnmarshal, err)
		}
	}
	return &logical.Response{
		Data: resp,
	}, nil
}
