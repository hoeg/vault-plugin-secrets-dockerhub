package dockerhub

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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
			b.configPaths(),
			b.tokenPaths(),
		),
		Secrets: []*framework.Secret{
			{
				Type:            "DockerHub",
				DefaultDuration: defaultTTL,
				Revoke:          b.handleRevokeToken,
				Fields: map[string]*framework.FieldSchema{
					tokenUsername: {
						Type:        framework.TypeString,
						Description: descTokenUsername,
					},
					tokenUUID: {
						Type:        framework.TypeString,
						Description: descTokenUUID,
					},
				},
			},
		},
	}

	return b, nil
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}
