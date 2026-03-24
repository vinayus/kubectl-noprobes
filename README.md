# kubectl-noprobes

A kubectl plugin that lists all pods and containers missing liveness or readiness probes across your cluster.

## Installation

### Manual

```bash
go install github.com/vinayus/kubectl-noprobes@latest
mv $(go env GOPATH)/bin/kubectl-noprobes /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/vinayus/kubectl-noprobes.git
cd kubectl-noprobes
go build -o kubectl-noprobes .
mv kubectl-noprobes /usr/local/bin/
```

## Usage

```bash
kubectl noprobes
```

### Example output

```
NAMESPACE    POD                          CONTAINER    MISSING
default      nginx-6d4cf56db6-xkq2p       nginx        liveness, readiness
kube-system  coredns-5dd5756b68-abc12     coredns      readiness
staging      api-server-7f9b4c6d5-zxp99   api          liveness
```

## Requirements

- `kubectl` configured with a valid kubeconfig
- Go 1.21+ (to build from source)
