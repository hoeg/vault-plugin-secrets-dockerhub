package token

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/config"
	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/dockerhub"
)

// pathPatternToken is the string used to define the base path of the config
// endpoint as well as the storage path of the config object.
var pathPatternToken = fmt.Sprintf("token/(%s)/(%s)", framework.GenericNameRegex(Username), framework.GenericNameRegex(Scope))

const (
	Username          = "username"
	DescTokenUsername = "Username for DockerHub."
	Scope             = "scope"
	DescTokenScope    = "Scope of the token."
	Label             = "label"
	DescTokenLabel    = "Name for the token to create."
	UUID              = "uuid"
	DescTokenUUID     = "The uuid for a generated token. Used for revokation."
)

const pathTokenHelpSyn = `

`

var pathTokenHelpDesc = "Issue an access token to Docker Hub with a given scope."

func Paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternToken,

			Fields: map[string]*framework.FieldSchema{
				Username: {
					Type:        framework.TypeString,
					Description: DescTokenUsername,
				},
				Scope: {
					Type:        framework.TypeString,
					Description: DescTokenScope,
				},
				Label: {
					Type:        framework.TypeString,
					Description: DescTokenLabel,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: handleCreate,
					Summary:  "Issue a new access token to Docker Hub.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: handleCreate,
					Summary:  "Issue a new access token to Docker Hub.",
				},
				logical.RevokeOperation: &framework.PathOperation{
					Callback: HandleRevoke,
					Summary:  "Revoke access token for Docker Hub.",
				},
			},
			HelpSynopsis:    pathTokenHelpSyn,
			HelpDescription: pathTokenHelpDesc,
			ExistenceCheck:  handleExistenceCheck,
		},
	}
}

func handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

func handleCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := data.Get(Username).(string)
	scope := data.Get(Scope).(string)
	l := data.Get(Label).(string)

	c, err := config.Retrieve(ctx, u, req.Storage)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(c.Scopes, scope) {
	 return nil, fmt.Errorf("requested scope %q is not permitted for user %q", scope, u)
	}

	dc := dockerhub.Client{
		Conf: c,
	}

	t, err := dc.NewToken(ctx, l, scope)
	if err != nil {
		return nil, err
	}

	logger := hclog.New(&hclog.LoggerOptions{})
	logger.Error(fmt.Sprintf("ttl: %s", c.TTL))
	return &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       c.TTL,
				Renewable: false,
			},
			InternalData: map[string]interface{}{
				"secret_type": "DockerHub",
				Username:      c.Username,
				Scope:         scope,
				UUID:          t.UUID,
			},
		},
		Data: map[string]interface{}{
			"token":  t.Token,
			UUID:     t.UUID,
			Username: c.Username,
			Scope:    scope,
		},
	}, nil
}

func HandleRevoke(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := data.Get(Username).(string)
	UUID := data.Get(UUID).(string)

	c, err := config.Retrieve(ctx, u, req.Storage)
	if err != nil {
		return nil, err
	}

	dc := dockerhub.Client{
		Conf: c,
	}

	err = dc.DeleteToken(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
