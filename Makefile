DOCKER_NAME := "foodlogiq-demo"
VERSION := "1.0.0-DEV"
PROJECT_NAME := "foodlogiq"

.PHONY: help
help:
	@echo 'Commands for running foodlogiq-demo:'
	@echo
	@echo 'Usage:'
	@echo '    make deps                Setup test dependencies (GinkGo and Mockery)'
	@echo '    make init                Init the project'
	@echo '    make build               Compile the project'
	@echo '    make run                 Run service'
	@echo '    make docker_build        Docker build'
	@echo '    make docker_run          Docker run'
	@echo '    make mocks               Generate mocks'
	@echo '    make test                Run tests on a compiled project'
	@echo

.PHONY: deps
deps: install-mockery install-ginkgo

.PHONY: install-ginkgo
install-ginkgo:
	go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo

.PHONY: install-mockery
install-mockery:
	go install github.com/vektra/mockery/v2@latest

.PHONY: init
init: tidy fmt vet lint

.PHONY: test
test: deps
	go test -cover `go list ./internal/... && go list ./pkg/...` -coverprofile=coverage.out

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo $(GOBUILD) ./cmd/$(PROJECT_NAME)/main.go

.PHONY: all
all: init test build

.PHONY: run
run:
	go run cmd/$(PROJECT_NAME)/main.go

.PHONY: docker_build
docker_build:
	docker build -t $(DOCKER_NAME):$(VERSION) .

.PHONY: docker_run
docker_run:
 	docker run -p 8000:8000 $(DOCKER_NAME)

mocks:
	mockery --dir pkg --all --output pkg/mocks

