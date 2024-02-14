# Vault Secrets Plugin for Docker Hub

[![CodeQL](https://github.com/hoeg/vault-plugin-secrets-dockerhub/actions/workflows/github-code-scanning/codeql/badge.svg?branch=master)](https://github.com/hoeg/vault-plugin-secrets-dockerhub/actions/workflows/github-code-scanning/codeql) [![Semgrep](https://github.com/hoeg/vault-plugin-secrets-dockerhub/actions/workflows/semgrep.yml/badge.svg?branch=master)](https://github.com/hoeg/vault-plugin-secrets-dockerhub/actions/workflows/semgrep.yml)

![rest](https://github.com/hoeg/vault-plugin-secrets-dockerhub/blob/master/pics/4wysdl.jpg)

Docker is used in many CI/CD piplines and accessing your private repositories should be made possible in a secure way. Using username and password for this is bad since these credentials have way to broad permissions. Access tokens on the other hand cannot change the password for an account and they can be restricted to specific namespaces thereby having a tighter scope than your username and password.

## Usage

To use the plugin you must rigster it. See the [Hashicorp Vault documentation](https://www.vaultproject.io/docs/commands/plugin/register) for the steps needed. The `Makefile` provides steps to test locally.

### Configure DockerHub account

First configure the credentials for the DockerHub account you want credentials from:

```bash
vault write dockerhub/config/$USERNAME password=$PASSWORD scopes=$SCOPE
```

where scopes is a comma separated list with the following valid values:`admin, write, read, public_read`.

`ttl` is optional. If it is not provided it will be set to the default `ttl` which is 5 minutes.

You can read the permissions using

```bash
vault read dockerhub/config/$USERNAME
```

The password will not be shown. Also it is not possible to update en existing configuration but a new one can be created. No validity checks are made when the config is written aside from validating the scopes.

### Creating tokens

Tokens issued by Vault will be revoked automatically after the `ttl` has expired. To issue a token run:

```bash
vault write dockerhub/token/$SCOPE label=$TOKEN_LABEL
```

By having scope as part of the path it is possible to restrict which scopes vault users are allowed to create credentials for.
