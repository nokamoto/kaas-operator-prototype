package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var controllers = []string{
	"pipeline",
	"kubernetescluster",
	"kubernetesclusterconfiguration",
}

var Default = All

// Install installs the necessary tools for the project.
func Install() error {
	tools := []string{
		"sigs.k8s.io/controller-tools/cmd/controller-gen",
		"golang.org/x/tools/cmd/goimports",
		"mvdan.cc/gofumpt",
		"github.com/google/ko",
		"sigs.k8s.io/controller-runtime/tools/setup-envtest",
		"github.com/bufbuild/buf/cmd/buf",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"connectrpc.com/connect/cmd/protoc-gen-connect-go",
		"github.com/google/yamlfmt/cmd/yamlfmt",
		"go.uber.org/mock/mockgen",
	}
	for _, tool := range tools {
		if err := sh.Run("go", "install", fmt.Sprintf("%s@latest", tool)); err != nil {
			return err
		}
	}
	return nil
}

// All runs all the necessary tasks to prepare the project.
func All() {
	mg.SerialDeps(
		Install,
		Build.Generate,
		Build.ControllerGenCRD,
		Build.ControllerGenObject,
		Build.ControllerGenRBAC,
		Build.ManagerYAML,
		Build.Buf,
		Build.Mock,
		Format.Proto,
		Format.Go,
		Format.YAML,
		Test,
		Build.InstallMCPServer,
	)
}
