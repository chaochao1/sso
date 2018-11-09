# Sso Service

This is the Sso service

Generated with

```
micro new github.com/cicdi-go/sso --namespace=go.micro --fqdn=go.micro.srv.sso --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.sso
- Type: srv
- Alias: sso

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./sso-srv
```

Build a docker image
```
make docker
```