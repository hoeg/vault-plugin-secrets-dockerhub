package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	loginEndpoint    = "https://hub.docker.com/v2/users/login"
	apiTokenEndpoint = "https://hub.docker.com/v2/api_tokens"
)

type Client struct{}

func NewClient(namespace string, s *logical.Storage) (*Client, error) {
	return nil, nil
}

type DockerHubToken struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

func (c Client) newToken(ctx context.Context, label string) (DockerHubToken, error) {
	apiToken, err := c.dockerHubAuth(ctx)
	if err != nil {
		return DockerHubToken{}, err
	}
	createToken := struct {
		TokenLabel string `json:"token_label"`
	}{
		TokenLabel: label,
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
	req.AddCookie(&http.Cookie{Name: "token", Value: apiToken})
	req.AddCookie(&http.Cookie{Name: "namespace", Value: c.Namespace})

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return DockerHubToken{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return DockerHubToken{}, fmt.Errorf("failed to createauth token: %d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
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

func (c Client) deleteToken(ctx context.Context, uuid string) error {
	apiToken, err := c.dockerHubAuth(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", apiTokenEndpoint, uuid), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "vault-docker-hub-secret")
	req.AddCookie(&http.Cookie{Name: "token", Value: apiToken})
	req.AddCookie(&http.Cookie{Name: "namespace", Value: c.Namespace})

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

func (c Client) dockerHubAuth(ctx context.Context) (string, error) {
	login := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: c.Username,
		Password: c.Password,
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
	b, err := ioutil.ReadAll(resp.Body)
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
