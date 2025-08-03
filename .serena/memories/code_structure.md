# Code Structure

- `cmd/`: Entrypoints (CLI, controllers, MCP server)
- `internal/`: Core logic (apiclient, cli, controller, infra, mcp, mock, service)
- `api/`: CRDs and proto definitions
- `pkg/`: Generated code and shared packages
- `magefiles/`: Mage automation scripts
- `config/`: CRD, manager, and RBAC YAMLs
- `docs/`: Design and documentation
- `hack/`: Boilerplate and templates
