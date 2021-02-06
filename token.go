package dockerhub

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathPatternToken is the string used to define the base path of the config
// endpoint as well as the storage path of the config object.
var pathPatternToken = fmt.Sprintf("token/(%s)/(%s)", framework.GenericNameRegex(tokenUsername), framework.GenericNameRegex(tokenNamespace))

const (
	fmtErrTokenMarshal = "failed to marshal token to JSON"
	fmtErrTokenPersist = "failed to persist token to storage"
	fmtErrTokenDelete  = "failed to delete token from storage"
)

const (
	tokenUsername      = "username"
	descTokenUsername  = "Username that has access to the namespace."
	tokenNamespace     = "namespace"
	descTokenNamespace = "Docker namespace to issue a token to."
	tokenLabel         = "label"
	descTokenLabel     = "Name for the token to create."
	tokenUuid          = "uuid"
	descTokenUuid      = "The uuid for a generated token. Used for revokation."
)

const pathTokenHelpSyn = `

`

var pathTokenHelpDesc = fmt.Sprintf(`
Issue an access token to Docker Hub for given namespace.`)

func (b *backend) tokenPaths() []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: pathPatternToken,

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

	c, err := b.Config(ctx, u, req.Storage)
	if err != nil {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}
	if ns != c.Namespace {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": "illegal namespace",
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

	return &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       c.TTL,
				Renewable: false,
				MaxTTL:    c.MaxTTL,
			},
			InternalData: map[string]interface{}{
				tokenUuid:     t.Uuid,
				tokenUsername: c.Username,
			},
		},
		Data: map[string]interface{}{
			"token": t.Token,
			"uuid":  t.Uuid,
		},
	}, nil
}

func (b *backend) handleRevokeToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	u := getStringFrom(data, tokenUsername)
	uuid := getStringFrom(data, tokenUuid)
	c, err := b.Config(ctx, u, req.Storage)
	if err != nil {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}
	err = c.DeleteToken(ctx, uuid)
	if err != nil {
		return &logical.Response{
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}
	return nil, nil
}
