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

	config := b.Config(u)
	if ns != config.Namespace {
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

	//store t
	tokenInfo := struct {
		Uuid  string
		Label string
	}{}
	ce, err := logical.StorageEntryJSON(tokenStoragePath(c.Username), c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfMarshal, err)
	}
	if err = req.Storage.Put(ctx, ce); err != nil {
		return nil, fmt.Errorf("%s: %w", fmtErrConfPersist, err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token": t.Token,
			"uuid":  t.Uuid,
		},
	}, nil
}

func tokenStoragePath(label) {
	return fmt.Sprintf("token/")
}

func (b *backend) handleRevokeToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) handleListToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
