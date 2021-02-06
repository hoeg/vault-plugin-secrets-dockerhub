package dockerhub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathPatternConfig is the string used to define the base path of the config
// endpoint as well as the storage path of the config object.
var pathPatternConfig = fmt.Sprintf("config/%s", framework.GenericNameRegex(configUsername))

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

const defaultTTL time.Duration = 5 * time.Minute

var pathConfigHelpDesc = fmt.Sprintf(``)

func (b *backend) configPaths() []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				configNamespace: {
					Type:        framework.TypeString,
					Description: descConfigNamespace,
					Required:    true,
				},
				configUsername: {
					Type:        framework.TypeString,
					Description: descConfigUsername,
					Required:    true,
				},
				configPassword: {
					Type:        framework.TypeString,
					Description: descConfigPassword,
					Required:    true,
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
					Callback: b.handleCreateConfig,
					Summary:  "Create a configuration for Docker Hub.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDeleteConfig,
					Summary:  "Deletes the secret at the specified location.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadConfig,
					Summary:  "Read the configuration for Docker Hub access tokens for a specific user.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCreateConfig,
					Summary:  "Update an existing configuration for Docker Hub.",
				},
			},
			HelpSynopsis:    pathConfigHelpSyn,
			HelpDescription: pathConfigHelpDesc,
			ExistenceCheck:  b.handleExistenceCheck,
		},
	}
}

//Config holds to values needed to issue a new Docker Hub access token
type Config struct {
	Namespace string        `json:"namespace"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	TTL       time.Duration `json:"ttl"`
	MaxTTL    time.Duration `json:"max_ttl"`
}

func (b *backend) Config(ctx context.Context, username, namespace string, s logical.Storage) (*Config, error) {
	c := &Config{}
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

	return c, nil
}

func (b *backend) handleCreateConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	c := Config{}
	c.Username = getStringFrom(data, configUsername)
	c.Password = getStringFrom(data, configPassword)
	c.Namespace = getStringFrom(data, configNamespace)
	c.TTL = defaultTTL

	fmt.Printf("Config is: %v/n", c)

	ce, err := logical.StorageEntryJSON(getStorePath(c.Username, c.Namespace), c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfMarshal, err)
	}
	if err = req.Storage.Put(ctx, ce); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfPersist, err)
	}
	return nil, nil
}

func (b *backend) handleDeleteConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := getStringFrom(data, configUsername)
	ns := getStringFrom(data, configNamespace)
	return nil, req.Storage.Delete(ctx, getStorePath(u, ns))
}

func getStorePath(u, ns string) string {
	return fmt.Sprintf("config/%s/%s", u, ns)
}

func (b *backend) handleReadConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := getStringFrom(data, configUsername)
	ns := getStringFrom(data, configNamespace)
	c, err := b.Config(ctx, u, ns, req.Storage)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}

	resp := make(map[string]interface{})

	if v := c.Username; v != "" {
		resp["username"] = v
	}
	if v := c.Password; v != "" {
		resp["password"] = v
	}
	if v := c.Namespace; v != "" {
		resp["namespace"] = v
	}
	resp["ttl"] = c.TTL

	return &logical.Response{
		Data: resp,
	}, nil
}

func (b *backend) handleListConfigs(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
