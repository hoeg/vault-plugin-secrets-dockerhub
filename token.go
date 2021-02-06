package dockerhub

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathPatternConfig is the string used to define the base path of the config
// endpoint as well as the storage path of the config object.
var pathPatternToken = fmt.Sprint("token/(%s)/(%s)", tokenUsername, tokenNamespace)

const (
	fmtErrTokenMarshal = "failed to marshal token to JSON"
	fmtErrTokenPersist = "failed to persist token to storage"
	fmtErrTokenDelete  = "failed to delete token from storage"
)

const (
	tokenUsername      = "username"
	descTokenUsername  = ""
	tokenNamespace     = "namespace"
	descTokenNamespace = "Docker namespace to issue a token to."
	tokenLabel         = "token-label"
	descTokenLabel     = "Name for the token to create."
)

const pathTokenHelpSyn = `

`

var pathTokenHelpDesc = fmt.Sprintf(`
Issue an access token to Docker Hub for given namespace.`)

func (b *backend) tokenPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				tokenUsername: {
					Type:        framework.TypeString,
					Description: descTokenUsername,
				},
				tokenNamespace: {
					Type:        framework.TypeString,
					Description: descTokenNamespace,
				},
				tokenLabel: {
					Type:        framework.TypeString,
					Description: descTokenLabel,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleCreateToken,
					Summary:  "Issue a new access token to Docker Hub.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCreateToken,
					Summary:  "Issue a new access token to Docker Hub.",
				},
				logical.RevokeOperation: &framework.PathOperation{
					Callback: b.handleRevokeToken,
					Summary:  "Revoke access token for Docker Hub.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleListToken,
					Summary:  "List issued access tokens for Docker Hub access tokens.",
				},
			},
			HelpSynopsis:    pathTokenHelpSyn,
			HelpDescription: pathTokenHelpDesc,
			ExistenceCheck:  b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleCreateToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := getStringFrom(data, tokenUsername)
	ns := getStringFrom(data, tokenNamespace)
	l := getStringFrom(data, tokenLabel)

	c, err := NewClient(u, ns, &req.Storage)
	if err != nil {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}

	t, err := c.NewToken(ctx, l)
	if err != nil {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}

	//store uuid + label

	return &logical.Response{
		Data: map[string]interface{}{
			"token": t.Token,
			"uuid":  t.Uuid,
		},
	}, nil
}

func (b *backend) handleRevokeToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) handleListToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
