.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: build

.PHONY: FORCE
FORCE:

.PHONY: build
build: ## builds binaries
	go build -o bin/gha ./cmd/gha

.PHONY: install
install: build ## installs the docker plugins
	cp bin/gha ~/bin/gha

.PHONY: check-line-endings
check-line-endings: ## check line endings
	! find . -name "*.go" -type f -exec file "{}" ";" | grep CRLF

.PHONY: fix-line-endings
fix-line-endings: ## fix line endings
	find . -type f -name "*.go" -exec sed -i -e "s/\r//g" {} +
