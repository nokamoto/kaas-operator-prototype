package main

import (
	"errors"
	"fmt"
	"slices"

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

var (
	configs = []string{
		"crd",
		"rbac",
		"manager",
	}
)

// Apply applies the Kubernetes manifests for local development.
func (Kind) Apply() error {
	for _, config := range configs {
		if err := sh.RunV("kubectl", "apply", "-f", fmt.Sprintf("./config/%s", config)); err != nil {
			return err
		}
	}
	return nil
}

// Clean removes the deployed resources from the Kind cluster.
func (Kind) Clean() error {
	// reverse the order of configs to ensure dependencies are cleaned up first
	reversedConfigs := make([]string, len(configs))
	copy(reversedConfigs, configs)
	slices.Reverse(reversedConfigs)
	var es []error
	for _, config := range reversedConfigs {
		if err := sh.RunV("kubectl", "delete", "-f", fmt.Sprintf("./config/%s", config)); err != nil {
			fmt.Printf("Error cleaning up %s: %v\n", config, err)
			es = append(es, err)
		}
	}
	return errors.Join(es...)
}
