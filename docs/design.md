
# Design Document

This document describes the design principles and architecture of the kaas-operator-prototype.

## 1. Overview

- Prototype for designing and experimenting with a KaaS (Kubernetes as a Service) controller
- Focus on flexibility, extensibility, and operational ease

## 2. Architecture

### Diagrams

#### CLI, gRPC Services and CRDs

```mermaid
flowchart LR
  cli[CLI]
  subgraph gRPC_Services["gRPC Services"]
    pipelineservice[PipelineService]
    k8sclusterservice[KubernetesClusterService]
  end
  subgraph CRDs["Custom Resource Definitions"]
    p[Pipeline]
    kc[KubernetesCluster]
    kccfg[KubernetesClusterConfiguration]
  end
  cli --> pipelineservice
  cli --> k8sclusterservice
  pipelineservice --> p
  k8sclusterservice --> kc
  k8sclusterservice --> kccfg
```

#### Controllers and CRDs

```mermaid
flowchart LR
  subgraph Controllers["Controllers"]
    pqc[PipelineQueueController]
    pc[PipelineController]
    kcc[KubernetesClusterController]
    kccc[KubernetesClusterConfigurationController]
  end
  subgraph CRDs["Custom Resource Definitions"]
    p[Pipeline]
    kc[KubernetesCluster]
    kccfg[KubernetesClusterConfiguration]
  end
  pqc --> p
  pc --> p
  kcc --> kc
  kccc --> kccfg
  kccc --> kc
```

### CLI

- Provides a command-line interface for users to interact with the system via gRPC services.

### gRPC Services

The system provides multiple gRPC services that enable external clients to interact with the Kubernetes Custom Resources (CRDs) managed by this project.

- **PipelineService**: Exposes create, read, and update operations for the Pipeline CRD.
- **KubernetesClusterService**: Exposes create, read, and update operations for the KubernetesCluster and KubernetesClusterConfiguration CRDs.
- This allows programmatic management and integration with other systems.

> [!NOTE]
> Delete operations may be restricted or handled separately depending on the use case and safety requirements.

### Kubernetes Controllers

- PipelineQueueController
- PipelineController
- KubernetesClusterController
- KubernetesClusterConfigurationController

### Custom Resource Definitions (CRDs)

- Pipeline
- KubernetesCluster
- KubernetesClusterConfiguration

---

*Details will be added as the project evolves.*
