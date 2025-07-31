package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Kind mg.Namespace

// Build builds the application for Kind.
func (Kind) Build() error {
	apps := []string{
		"./cmd/pipelinecontroller",
	}
	env := map[string]string{
		"KO_DOCKER_REPO": "kind.local",
	}
	for _, app := range apps {
		if err := sh.RunWithV(env, "ko", "build", "--base-import-paths", app); err != nil {
			return err
		}
	}
	return nil
}

// Apply applies the Kubernetes manifests for local development.
func (Kind) Apply() error {
	return sh.RunV("kubectl", "apply", "-f", "./config/deployment")
}

// Clean removes the deployed resources from the Kind cluster.
func (Kind) Clean() error {
	return sh.RunV("kubectl", "delete", "-f", "./config/deployment")
}
