# kaas-operator-prototype

Testing ground for KaaS controller architecture

## Documentation

- [Design Document](docs/design.md)

## Development Workflow (Mage)

This project uses [Mage](https://magefile.org/) to automate development tasks.

See the full list of available tasks with:

```sh
mage -l
```

### Common Use Cases

#### All-in-one build

Run all necessary build steps for the project.

```sh
mage
```

#### Local development with Kind

Build and deploy to a local Kind cluster:

```sh
mage kind:build
mage kind:apply
```

Remove deployed applications from Kind:

```sh
mage kind:clean
```

#### Setup for Serena MCP

Prepare the environment for Serena MCP development:

```sh
direnv allow
mage python:venv
mage python:install
```
