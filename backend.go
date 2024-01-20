package dockerhub

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/config"
	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/token"
)

const dockerHubHelp = `
The DockerHub secrets backend will create a temporary access token for Docker Hub.
`

// backend wraps the backend framework
type backend struct {
	*framework.Backend
	configLock *sync.Mutex
}

var _ logical.Factory = Factory

// Factory configures and returns DockerHub backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	b, err := newBackend()
	if err != nil {
		return nil, err
	}
	b.Logger().Info("plugin backend initialization started")

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.Logger().Info("plugin backend successfully initialised")
	return b, nil
}

func newBackend() (*backend, error) {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(dockerHubHelp),
		BackendType: logical.TypeLogical,
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/",
			},
		},
		Paths: framework.PathAppend(
			config.Paths(),
			token.Paths(),
		),
		Secrets: []*framework.Secret{
			{
				Type:            "DockerHub",
				DefaultDuration: config.DefaultTTL,
				Revoke:          token.HandleRevoke,
				Fields: map[string]*framework.FieldSchema{
					token.Username: {
						Type:        framework.TypeString,
						Description: token.DescTokenUsername,
					},
					token.Namespace: {
						Type:        framework.TypeString,
						Description: token.DescTokenNamespace,
					},
					token.UUID: {
						Type:        framework.TypeString,
						Description: token.DescTokenUUID,
					},
				},
			},
		},
	}

	return b, nil
}
