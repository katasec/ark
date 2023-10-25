.ONESHELL:
SHELL = /bin/bash
.DEFAULT_GOAL := help
.PHONY: worker server checkbuild workerpush


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
worker: ## Build Worker container
	@IMAGE_NAME=arkworker;
	@source ./scripts/build.sh;
	@build
workerpush: ## Build and push worker container
	@IMAGE_NAME=arkworker;
	@source ./scripts/build.sh;
	@buildAndPush
checkbuild: ## Check app can build
	go build -o /dev/null
server: ## Build Server container
	@IMAGE_NAME=arkserver;
	@source ./scripts/build.sh;
	@build
#	docker build -t ghcr.io/katasec/arkapi:0.01 .
rserver: ## Run Server
	@go install
	@ark server
rworker: ## Run Server
	@go install
	@ark worker start