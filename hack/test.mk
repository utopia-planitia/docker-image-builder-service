
.PHONY: test
test: ##@testing Runs all tests.
	@$(MAKE) go-test
	@$(MAKE) queue-test

.PHONY: go-test
go-test: ##@testing Runs go (unit) tests.
	go test -race ./...

.PHONY: queue-test
queue-test: .devtools ##@testing Runs more parallel builds then workers
	./hack/queue-test.sh
