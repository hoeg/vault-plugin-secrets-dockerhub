# Vault Secrets Plugin for Docker Hub

![rest](https://github.com/hoeg/vault-plugin-secrets-dockerhub/blob/master/pics/4wysdl.jpg)

Docker is used in many CI/CD piplines and accessing your private repositories should be made possible in a secure way. Using username and password for this is bad since these credentials have way to broad permissions.

## Usage

### Register the plugin

### Configure DockerHub account

`vault write dockerhub/config/<username> password=<password> namespace=<namespace>`

`ttl` is optional. If it is not provided it will be set to the default `ttl` which is 5 minutes.

### Creating tokens

`vault write dockerhub/token/<username>/<namespace> label=<token label>`
