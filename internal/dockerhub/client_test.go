package dockerhub

import (
	"context"
	"testing"

	"github.com/hoeg/vault-plugin-secrets-dockerhub/internal/value"
)

func TestCreateToken(t *testing.T) {
	t.Skip()
	conf := value.Config{
		Scopes:   []string{"hoeg"},
		Username: "hoeg",
		Password: "",
	}

	client := Client{
		Conf: &conf,
	}

	token, err := client.NewToken(context.Background(), "tes1", "hoeg")
	if err != nil {
		t.Fatal(err)
	}
	if token.UUID == "" {
		t.Fatal("no valid token returned")
	}
}

func TestDeleteToken(t *testing.T) {
	t.Skip()
	conf := value.Config{
		Scopes:   []string{"hoeg"},
		Username: "hoeg",
		Password: "",
	}

	client := Client{
		Conf: &conf,
	}

	err := client.DeleteToken(context.Background(), "c071fa97-fe94-41d5-b1f8-96af5604e3d3")
	if err != nil {
		t.Fatal(err)
	}
}
