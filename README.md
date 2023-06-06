# Donna

![Logo](./docs/logo-readme.png)

Minimal personal CRM.

[![hydrun CI](https://github.com/pojntfx/donna/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/donna/actions/workflows/hydrun.yaml)
[![Docker CI](https://github.com/pojntfx/donna/actions/workflows/docker.yaml/badge.svg)](https://github.com/pojntfx/donna/actions/workflows/docker.yaml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/donna.svg)](https://pkg.go.dev/github.com/pojntfx/donna)
[![Matrix](https://img.shields.io/matrix/donnadev:matrix.org)](https://matrix.to/#/#donnadev:matrix.org?via=matrix.org)
[![Binary Downloads](https://img.shields.io/github/downloads/pojntfx/donna/total?label=binary%20downloads)](https://github.com/pojntfx/donna/releases)

## Overview

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

## Installation

### Hosted Demo

If you just want to quickly evaluate Donna, check out the hosted demo:

[<img src="https://github.com/pojntfx/webnetesctl/raw/main/img/launch.png" width="240">](https://donna-demo.vercel.app/)

### Containerized

You can get the OCI image like so:

```shell
$ podman pull ghcr.io/pojntfx/donna
```

### Natively

Static binaries are available on [GitHub releases](https://github.com/pojntfx/donna/releases).

On Linux, you can install them like so:

```shell
$ curl -L -o /tmp/donna "https://github.com/pojntfx/donna/releases/latest/download/donna.linux-$(uname -m)"
$ sudo install /tmp/donna /usr/local/bin
```

On macOS, you can use the following:

```shell
$ curl -L -o /tmp/donna "https://github.com/pojntfx/donna/releases/latest/download/donna.darwin-$(uname -m)"
```

On Windows, the following should work (using PowerShell as administrator):

```shell
PS> Invoke-WebRequest https://github.com/pojntfx/donna/releases/latest/download/donna.windows-x86_64.exe -OutFile \Windows\System32\donna.exe
```

You can find binaries for more operating systems and architectures on [GitHub releases](https://github.com/pojntfx/donna/releases).

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build and start a development version of donna locally, run the following:

```shell
$ git clone https://github.com/pojntfx/donna.git
$ cd donna
$ make depend
$ docker rm -f donna-postgres && docker run -d --name donna-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=donna postgres
$ docker exec donna-postgres bash -c 'until pg_isready; do sleep 1; done'
$ export OIDC_ISSUER='https://pojntfx.eu.auth0.com/' OIDC_CLIENT_ID='dyMxiRh1v2o8ALJcxN1WiHbmRygqNyno' OIDC_REDIRECT_URL='http://localhost:1337/authorize'
$ go run ./cmd/donna
```

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#donnadev:matrix.org?via=matrix.org)!

## License

Donna (c) 2023 Felicitas Pojtinger, Daniel Hiller and contributors

SPDX-License-Identifier: AGPL-3.0
