package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	dockerhub "github.com/hashicorp/vault-guides/plugins/vault-plugin-secrets-mock"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{})
	logger.Info("started plugin")
	apiClientMeta := &api.PluginAPIClientMeta{}

	flags := apiClientMeta.FlagSet()
	if err := flags.Parse(os.Args[1:]); err != nil {
		logger.Error("Failed to parse flags", err)
		os.Exit(1)
	}

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: dockerhub.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}


vault write sys/plugins/catalog/vault-plugin-secrets-dockerhub sha_256=$( shasum -a 256 ./vault/plugins/vault-plugin-secrets-dockerhub | cut -d " " -f1) command="vault-plugin-secrets-dockerhub"