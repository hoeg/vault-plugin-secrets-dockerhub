# Vault Secrets Plugin for Docker Hub

![rest](https://github.com/hoeg/vault-plugin-secrets-dockerhub/blob/master/pics/4wysdl.jpg)

Docker is used in many CI/CD piplines and accessing your private repositories should be made possible in a secure way. Using username and password for this is bad since these credentials have way to broad permissions. Access tokens on the other hand cannot change the password for an account and they can be restricted to specific namespaces thereby having a tighter scope than your username and password.

## Usage

To use the plugin you must rigster it. See the [Hashicorp Vault documentation](https://www.vaultproject.io/docs/commands/plugin/register) for the steps needed. The `Makefile` provides steps to test locally.

### Configure DockerHub account

First configure the credentials for the DockerHub account you want credentials from:

```
vault write dockerhub/config/<username> password=<password> namespace=<namespace>
```

 where namespace is a comma separated list of namespaces.

`ttl` is optional. If it is not provided it will be set to the default `ttl` which is 5 minutes.

You can read the permissions using

```
vault read dockerhub/config/<username>
```

 The password will not be shown. Also it is not possible to update en existing configuration but a new one can be created. No validity checks are made when the config is written.

### Creating tokens

Tokens issued by Vault will be revoked automatically after the `ttl` has expired. To issue a token run:

```
vault write dockerhub/token/<username>/<namespace> label=<token label>
```

By having namespace as part of the path it is possible to restrict which namespace vault users are allowed to create credentials for.


## Disclaimer

This plugin is build as an educational exercise in a day to learn about the Hashicorp Vault plugin structure. No garuantees are made about its security or stability (see the lack of tests). Use at your own risk...


## TODO

- List configurations
- A lot of cleanup!!