
.PHONY: lint
lint: ##@linting Runs all linters.
	@$(MAKE) lint-bash
	@$(MAKE) lint-go

.PHONY: lint-bash
lint-bash: ##@linting Lint Bash scripts.
	shellcheck etc/*.sh caching-scripts/*.sh

.PHONY: lint-go
lint-go: ##@linting Lint Go code.
	gometalinter ./...
