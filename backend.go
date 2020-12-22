package dockerhub

import (
	"context"
	"fmt"
	"strings"

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
}

var _ logical.Factory = Factory

// Factory configures and returns DockerHub backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b.Logger().Info("plugin backend initialization started")
	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

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
		Paths:       framework.PathAppend(
		//	b.tokenPaths(),
		//b.configPaths(),
		),
		Invalidate: b.Invalidate,
	}

	return b, nil
}

// Invalidate resets the plugin. It is called when a key is updated via
// replication.
func (b *backend) Invalidate(_ context.Context, key string) {
	if key == pathPatternConfig {
		// Configuration has changed so reset the client.
		//b.clientLock.Lock()
		//b.client = nil
		//b.clientLock.Unlock()
	}
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}
