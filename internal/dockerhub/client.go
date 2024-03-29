package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/value"
)

const (
	loginEndpoint    = "https://hub.docker.com/v2/users/login"
	apiTokenEndpoint = "https://hub.docker.com/v2/access-tokens"
)

type DockerHubToken struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
}

type Client struct {
	Conf *value.Config
}

// NewToken creates new access token and stores the uuid together with the label for lookup.
func (c *Client) NewToken(ctx context.Context, label, scope string) (DockerHubToken, error) {
	apiToken, err := c.dockerHubAuth(ctx)
	if err != nil {
		return DockerHubToken{}, err
	}
	createToken := struct {
		TokenLabel string   `json:"token_label"`
		Scopes     []string `json:"scopes"`
	}{
		TokenLabel: label,
		Scopes:     []string{fmt.Sprintf("repo:%s", scope)},
	}
	payload, err := json.Marshal(createToken)
	if err != nil {
		return DockerHubToken{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiTokenEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return DockerHubToken{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "vault-docker-hub-secret")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return DockerHubToken{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return DockerHubToken{}, fmt.Errorf("failed to create auth token: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return DockerHubToken{}, err
	}
	t := DockerHubToken{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return DockerHubToken{}, err
	}
	return t, nil
}

// DeleteToken will delete at token that is associated with the uuid.
func (c *Client) DeleteToken(ctx context.Context, UUID string) error {
	apiToken, err := c.dockerHubAuth(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", apiTokenEndpoint, UUID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "vault-docker-hub-secret")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete auth token %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) dockerHubAuth(ctx context.Context) (string, error) {
	login := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: c.Conf.Username,
		Password: c.Conf.Password,
	}
	payload, err := json.Marshal(login)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	t := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return "", err
	}
	return t.Token, nil
}
