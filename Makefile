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
	@echo '    make docker-build        Docker build'
	@echo '    make docker-run          Docker run'
	@echo '    make mocks               Generate mocks'
	@echo '    make test                Run tests on a compiled project'
	@echo '    make compose             docker-compose build and up in a detached state'
	@echo '    make run-demo        	compose the app, and then run the 4 test REST requests'
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
	go test -cover `go list ./pkg/...` -coverprofile=coverage.out

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo $(GOBUILD) ./cmd/$(PROJECT_NAME)/main.go

.PHONY: all
all: init test build

.PHONY: run
run:
	go run cmd/$(PROJECT_NAME)/main.go

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_NAME):$(VERSION) .

.PHONY: docker-run
docker-run:
 	docker run -p 8000:8000 $(DOCKER_NAME)

.PHONY: mocks
mocks:
	mockery --dir pkg --all --output pkg/mocks

.PHONY: compose
compose:
	docker-compose build && docker-compose up --detach

.PHONY: run-demo
run-demo: compose
	go run cmd/test/main.go

