# kaas-operator-prototype

Testing ground for KaaS controller architecture

## Documentation

- [Design Document](docs/design.md)

## Development Workflow (Mage)

This project uses [Mage](https://magefile.org/) to automate development tasks.

### Install dependencies

To install Go dependencies for development, run:

```sh
mage install
```

### Build

To generate Custom Resource Definitions (CRDs) using controller-gen, run:

```sh
mage build:controllerGenCRD
```

### Formatting

To format Go code and tidy dependencies, run:

```sh
mage format:go
```
