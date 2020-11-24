# CLI

The Krane [CLI](https://github.com/krane/cli) allows you to interact with Krane to create container resources.

## Installing

```
npm i -g @krane/cli
```

## Authenticating

Krane uses [private and public key authentication](https://en.wikipedia.org/wiki/Public-key_cryptography). Both keys are used for ensuring authenticity of incoming request.

1. Create the public and private key 

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -f $HOME/.ssh/krane -m 'PEM'

-t type
-b bytes
-C comments
-m key format
-f output file
```

This will generate 2 different keys, a `private` & `public (.pub)` key.

2. 

Place the public key on the host machine where Krane is running appended to `~/.ssh/authorized_keys`.

The `private key` is kept on the user's machine, the `public key` is stored where Krane is running and appended to `~/.ssh/authorized_keys`

Now you can try logging in. The CLI will prompt you to select the public key your just created. This will be used for authenticating with the private key from the Krane server.

```
krane login
```

## Commands

### login

Authenticate with Krane

```
krane login
```

### config

List the deployment configuration

```
krane config <deployment>
```

### delete

Delete a deployment

```
krane delete <deployment>
```

### describe

Describe a deployment. This provides details on the containers part of the deployment.

```
krane describe <deployment>
```

### deploy

Deploy or update an application.

```
krane deploy -f </path/to/krane.json>
```

### list

List all deployments

```
krane list
```

### secrets

List all deployment secrets

```
krane secrets list <deployment>
```

Add a deployment secret

```
krane secrets add <deployment> -k <key> -v <value>
```

Delete a deployment secret

```
krane secrets delete <deployment> -k <key>
```