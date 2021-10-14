# env
export GO111MODULE=on
export CGO_ENABLED=0

# project metadta
NAME         := kube-throttler
VERSION      ?= $(if $(RELEASE),$(shell $(GIT_SEMV) now),$(shell $(GIT_SEMV) patch -p))
REVISION     := $(shell git rev-parse --short HEAD)
IMAGE_PREFIX ?= 
IMAGE_TAG    ?= $(if $(RELEASE),$(VERSION),$(VERSION)-$(REVISION))
LDFLAGS      := -ldflags="-s -w -X \"github.com/everpeace/kube-throttler/cmd.Version=$(VERSION)\" -X \"github.com/everpeace/kube-throttler/cmd.Revision=$(REVISION)\" -extldflags \"-static\""
OUTDIR       ?= ./dist

.DEFAULT_GOAL := build

.PHONY: fmt
fmt:
	$(GO_IMPORTS) -w cmd/ pkg/
	$(GO_LICENSER) --licensor "Shingo Omura"

.PHONY: lint
lint: fmt
	$(GOLANGCI_LINT) run --config .golangci.yml --deadline 30m

.PHONY: build
build: fmt lint
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o $(OUTDIR)/$(NAME) .

.PHONY: install
install:
	kubectl apply -f ./deploy/crd.yaml

.PHONY: generate
generate: codegen crd

.PHONY: codegen
codegen:
	./hack/update-codegen.sh
	$(GO_LICENSER) --licensor "Shingo Omura"

.PHONY: crd
crd:
	$(CONTROLLER_GEN) crd paths=./pkg/apis/... output:stdout > ./deploy/crd.yaml

.PHONY: build-only
build-only: 
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o $(OUTDIR)/$(NAME) .

.PHONY: test
test: fmt lint
	go test  ./...

.PHONY: clean
clean:
	rm -rf "$(OUTDIR)"

.PHONY: build-image
build-image:
	docker build -t $(shell make -e docker-tag) --build-arg RELEASE=$(RELEASE) --build-arg VERSION=$(VERSION) --target runtime .
	docker tag $(shell make -e docker-tag) $(IMAGE_PREFIX)$(NAME):$(VERSION)  # without revision

.PHONY: push-image
push-image:
	docker push $(shell make -e docker-tag)
	# without revision
	docker push $(IMAGE_PREFIX)$(NAME):$(VERSION)
	# latest (update only in release)
	$(if $(RELEASE), docker tag $(shell make -e docker-tag) $(IMAGE_PREFIX)$(NAME):latest)
	$(if $(RELEASE), docker push $(IMAGE_PREFIX)$(NAME):latest)  

.PHONY: docker-tag
docker-tag:
	@echo $(IMAGE_PREFIX)$(NAME):$(IMAGE_TAG)

#
# Release
#
guard-%:
	@ if [ "${${*}}" = "" ]; then \
    echo "Environment variable $* is not set"; \
		exit 1; \
	fi
.PHONY: release
release: guard-RELEASE guard-RELEASE_TAG
	git diff --quiet HEAD || (echo "your current branch is dirty" && exit 1)
	git tag $(RELEASE_TAG) $(REVISION)
	git push origin $(RELEASE_TAG)


#
# dev setup
#
.PHONY: setup
DEV_TOOL_PREFIX = $(shell pwd)/.dev
GIT_SEMV = $(DEV_TOOL_PREFIX)/bin/git-semv
GOLANGCI_LINT = $(DEV_TOOL_PREFIX)/bin/golangci-lint
GO_LICENSER = $(DEV_TOOL_PREFIX)/bin/go-licenser 
GO_IMPORTS = $(DEV_TOOL_PREFIX)/bin/goimports
CONTROLLER_GEN = $(DEV_TOOL_PREFIX)/bin/controller-gen
setup:
	$(call go-get-tool,$(GO_IMPORTS),golang.org/x/tools/cmd/goimports)
	$(call go-get-tool,$(GO_LICENSER),github.com/elastic/go-licenser)
	$(call go-get-tool,$(GIT_SEMV),github.com/linyows/git-semv/cmd/git-semv)
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.1)
	cd $(shell go env GOPATH) && \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(DEV_TOOL_PREFIX)/bin v1.27.0

# go-get-tool will 'go get' any package $2 and install it to $1.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(DEV_TOOL_PREFIX)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

#
# local development
# TIPS: You can change loglevel dynamicaly:
#   $ curl curl -XPUT --data 'N' localhost:10251/debug/flags/v
#
KUBECONFIG ?= $(HOME)/.kube/config
THROTTLER_NAME ?= kube-throttler
SCHEDULER_NAME ?= my-scheduler
.PHONY: dev-scheduler-conf
dev-scheduler-conf:
	mkdir -p .dev
	KUBECONFIG=$(KUBECONFIG) \
	THROTTLER_NAME=$(THROTTLER_NAME) \
	SCHEDULER_NAME=$(SCHEDULER_NAME) \
	envsubst < ./hack/dev/scheduler-config.yaml.template > ./hack/dev/scheduler-config.yaml

.PHONY: dev-run
dev-run: dev-scheduler-conf
	go run main.go kube-scheduler \
		--config=./hack/dev/scheduler-config.yaml \
		-v=3

.PHONY: dev-run-debug
dev-run-debug: dev-scheduler-conf
	dlv debug --headless --listen=0.0.0.0:2345 --api-version=2 --log main.go -- kube-scheduler \
		--config=./hack/dev/scheduler-config.yaml \
		--kubeconfig=$(HOME)/.kube/config \
		--v=3