package dockerhub

import (
	"context"
	"testing"

	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/config"
)

func TestCreateToken(t *testing.T) {
	t.Skip()
	conf := config.Config{
		Namespace: []string{"hoeg"},
		Username:  "hoeg",
		Password:  "",
	}

	client := Client{
		Conf: &conf,
	}

	token, err := client.NewToken(context.Background(), "test", "hoeg")
	if err != nil {
		t.Fatal(err)
	}
	if token.UUID == "" {
		t.Fatal("no valid token returned")
	}
}
