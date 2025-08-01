package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nokamoto/kaas-operator-prototype/internal/service/cluster"
	"github.com/nokamoto/kaas-operator-prototype/internal/service/longrunningoperation"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
)

func main() {
	mux := http.NewServeMux()
	path, handler := v1alpha1connect.NewClusterServiceHandler(cluster.ClusterService{})
	mux.Handle(path, handler)
	path, handler = v1alpha1connect.NewLongRunningOperationServiceHandler(longrunningoperation.LongRunningOperationService{})
	mux.Handle(path, handler)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	slog.Info("starting server", "address", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
