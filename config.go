package dockerhub

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathPatternConfig is the string used to define the base path of the config
// endpoint as well as the storage path of the config object.
const pathPatternConfig = "config"

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

var pathConfigHelpDesc = fmt.Sprintf(``)

func (b *backend) configPaths() []*framework.Path {
	return []*framework.Path{
		{
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
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDeleteConfig,
					Summary:  "Deletes the secret at the specified location.",
				},
				logical.RevokeOperation: &framework.PathOperation{
					Callback: b.handleRevokeConfig,
					Summary:  "Revoke access token for Docker Hub.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleListConfig,
					Summary:  "List configurations for Docker Hub access tokens.",
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
	Namespace string
	Username  string
	Password  string
}

func (b *backend) handleCreateConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	c, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	ns := data.Get(configNamespace)
	if s, ok := ns.(string); ok && s != "" {
		c.Namespace = s
	}
	u := data.Get(configUsername)
	if s, ok := u.(string); ok && s != "" {
		c.Username = s
	}
	p := data.Get(configPassword)
	if s, ok := p.(string); ok && s != "" {
		c.Password = s
	}

	fmt.Printf("Config is: %v/n", c)

	ce, err := logical.StorageEntryJSON(pathPatternConfig, c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfMarshal, err)
	}
	if err = req.Storage.Put(ctx, ce); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfPersist, err)
	}
	return nil, nil
}

func (b *backend) Config(ctx context.Context, s logical.Storage) (*Config, error) {
	c := &Config{}

	entry, err := s.Get(ctx, pathPatternConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfRetrieval, err)
	}

	if entry == nil || len(entry.Value) == 0 {
		return c, nil
	}

	if err := entry.DecodeJSON(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfUnmarshal, err)
	}

	return c, nil
}

func (b *backend) handleDeleteConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) handleRevokeConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) handleListConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
