
BATS_TESTS=tests

.PHONY: test
test: ##@testing Runs all tests.
	@$(MAKE) go-test
	@$(MAKE) end-to-end-test
	@$(MAKE) queue-test

.PHONY: go-test
go-test: ##@testing Runs go (unit) tests.
	go test -race ./...

.PHONY: end-to-end-test
end-to-end-test: .devtools ##@testing Runs end to end tests.
	@docker run -ti --rm \
		--dns 10.96.0.10 --dns-search container-image-builder.svc.cluster.local \
		-e DOCKER_HOST=tcp://docker:2375 \
		-v $(PWD):/project -w /project \
		utopiaplanitia/docker-image-builder-devtools:latest bats ${BATS_TESTS}

.PHONY: queue-test
queue-test: .devtools ##@testing Runs more parallel builds then workers
	./hack/queue-test.sh
