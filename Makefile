DEV_MAKEFILES ?= tools/dev/makefiles
KUBE_NAMESPACE ?= todo

include $(DEV_MAKEFILES)/go.mk
include $(DEV_MAKEFILES)/go-sqlc.mk
include $(DEV_MAKEFILES)/openapi.mk
include $(DEV_MAKEFILES)/docker.mk
include $(DEV_MAKEFILES)/kind.mk
include $(DEV_MAKEFILES)/skaffold.mk

build: build-goose build-todo-service

bootstrap: bootstrap-deployment

.PHONY: build-goose build-todo-service bootstrap-deployment

build-goose: ## Build goose
	$(info $(_bullet) Building <goose>)
	$(GO) build -o bin/goose ./cmd/goose

build-todo-service: ## Build todo-service
	$(info $(_bullet) Building <todo-service>)
	$(GO) build -o bin/todo-service ./services/todo

bootstrap-deployment: $(KUBECTL) ## Bootstrap deployment
	$(info $(_bullet) Bootstraping <deployment>)
	$(KUBECTL) apply --context $(BOOTSTRAP_CONTEXT) -k ops/bootstrap/overlays/local
