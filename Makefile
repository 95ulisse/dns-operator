# Image URL to use all building/pushing image targets
VERSION ?= 0.1
IMG ?= 95ulisse/dns-operator:$(VERSION)
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: dns-operator

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build dns-operator binary
dns-operator: generate fmt vet
	go build -o bin/dns-operator main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install-crd: manifests
	kubectl apply -f config/crd

# Uninstall CRDs from a cluster
uninstall-crd: manifests
	kubectl delete -f config/crd

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=dns-operator-role webhook paths="./..." output:crd:artifacts:config=config/crd
	mkdir -p config/release/v$(VERSION)
	( for i in config/{crd,rbac,deployment}/*.yaml; do cat "$$i"; echo '---'; done ) > config/release/v$(VERSION)/all-in-one.yaml

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/release/v$(VERSION)/all-in-one.yaml

# Un-deploy controller in the configured Kubernetes cluster in ~/.kube/config
undeploy: manifests
	kubectl delete -f config/release/v$(VERSION)/all-in-one.yaml

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet -structtag=false ./...

# Generate code
generate: controller-gen
	go generate ./...
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
