.ONESHELL:
SHELL = /bin/bash
.DEFAULT_GOAL := help
.PHONY: worker server checkbuild workerpush


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
worker: ## Build Worker container
	@source ./scripts/build.sh && IMAGE_NAME=arkworker && build
workerpush: ## Build and push worker container
	@source ./scripts/build.sh && IMAGE_NAME=arkworker && buildAndPush
checkbuild: ## Check app can build
	@go build -o /dev/null
server: ## Build Server container
	@source ./scripts/build.sh && IMAGE_NAME=arkserver && build
#	docker build -t ghcr.io/katasec/arkapi:0.01 .
runserver: ## Run Server
	@go install
	@ark server
runworker: ## Run Worker
	@go install
	@ark worker start

publishazure: ## Publish to Azure functions
	@source ./scripts/build.sh && publishAzure	

runlocal: ## Run Azure func locally
	@source ./scripts/build.sh && runlocal	