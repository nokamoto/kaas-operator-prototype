# What to do when a task is completed

- Format code: `mage format:go`, `mage format:proto`, `mage format:yaml`
- Run tests: `mage test`
- Ensure CRDs and RBAC are up to date: `mage build:controllerGenCRD`, `mage build:controllerGenRBAC`
- Update dependencies if needed: `mage install`
- Commit and push changes via git
- Follow standard code review and PR process if contributing
