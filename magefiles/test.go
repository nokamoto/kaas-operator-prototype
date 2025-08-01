package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Test runs the tests for the project.
func Test() error {
	s, err := sh.Output("setup-envtest", "use", "-p", "path")
	if err != nil {
		return fmt.Errorf("failed to set up envtest: %w", err)
	}
	fmt.Printf("Using envtest at: %s\n", s)
	env := map[string]string{
		"KUBEBUILDER_ASSETS": s,
	}
	return sh.RunWithV(env, "go", "test", "-v", "./...")
}
