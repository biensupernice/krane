# Krane Config

Creating a deployment using Krane starts with a **single file**, this file contains the deployment configuration used when creating container resources. The deployment configuration can be stored at the root of your project, a ci repository or directory on your local machine. When using the [CLI](cli) you can reference the location of this deployment configuration so location. The simplest configuration file could look like:

```json
{
  "name": "hello-world-app",
  "image": "hello-world",
  "alias": ["hello.localhost"]
}
```

The above **deployment configuration** sets up:

1. A deployment called **hello-world-app**
2. A container running the image **hello-world**
3. A redirect alias **hello.localhost** automatically setup by Krane

See all deployment configuration [options](deployment-configuration?id=options)

## Options

The different properties you can specificy in a deployment configuration file.

> A common pattern is to have a `krane.json` at the root of your project

### name

The name of your deployment.

- required: `yes`

### registry

The container registry to use when pulling images.

- required: `false`
- default: `docker.io`

### image

The image used when pulling, creating and running your deployments containers.

- required: `true`

### ports

The ports to expose from the container to the host machine.

> ⚠️ Ports are discouraged since port conflicts can become an uwanted side-effect. Instead, Krane uses a reverse proxy that handles exposing your containers using [`aliases`](deployment-configuration?id=alias)

- required: `false`

```json
{
  "ports": {
    "80": "8080"
  }
}
```

### env

The environment variables passed to the containers part of a deployment.

> ⚠️ Environment variables should not contain any sensitive data, use [`secrets`](deployment-configuration?id=secrets) instead.

- required: `false`

```json
{
  "env": {
    "NODE_ENV": "dev",
    "PORT": "8080"
  }
}
```

### secrets

Secrets are used when you want to pass sensitive information to your deployments.

> You can add deployment secrets using the krane [`cli`](cli?id=secrets)

- required: `false`

```json
{
  "secrets": {
    "SECRET_TOKEN": "@MY_SECRET_TOKEN",
    "PROXY_API_URL": "@SOME_PROXY_API_URL"
  }
}
```

### tag

The tag used when pulling the image.

- required: `false`
- default: `latest`

### volumes

The volumes to mount from the container to the host.

- required: `false`

```json
{
  "volumes": {
    "/host/path": "/container/path"
  }
}
```

### alias

Ingress alias for your deployment 

required: `false`

```json
{
  "alias": [
    "app2.example.com",
    "app2-dev.example.com",
    "app2-mybranch.example.com"
  ]
}
```

### command

A custom command to start the containers.

- required: `false`

```json
{
  "command": "npm run start --prod"
}
```