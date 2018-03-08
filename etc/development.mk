
.PHONY: cli
cli: ##@development Opens a command line interface with development tools.
	docker build -f docker/devtools/Dockerfile  -t utopiaplanitia/docker-image-builder-devtools:latest .
	docker run -ti --rm \
		--dns 10.96.0.10 --dns-search container-image-builder.svc.cluster.local \
		-e DOCKER_HOST=tcp://docker:2375 \
		-e CACHE_ENDPOINT_SERVER=minio \
		-e CACHE_ENDPOINT_PORT=9000 \
		-e CACHE_BUCKET=image-layers \
		-e CACHE_ACCESS_KEY=8Q9U4RBHKKB6HU70SRZ1 \
		-e CACHE_SECRET_KEY=oxxT2iqBlW6lgaDVe8ll6mP8z/OSVIUnn9cB4+Q0 \
		-v $(PWD):/project -w /project \
		utopiaplanitia/docker-image-builder-devtools:latest sh

.PHONY: test
test: ##@testing Runs all tests.
	@$(MAKE) go-test
	@$(MAKE) end-to-end-test
	@$(MAKE) queue-test

.PHONY: go-test
go-test: ##@testing Runs go (unit) tests.
	go test -race ./...

.PHONY: end-to-end-test
end-to-end-test: ##@testing Runs end to end tests.
	@docker build -q -f docker/devtools/Dockerfile -t utopiaplanitia/docker-image-builder-devtools:latest .
	@docker run -ti --rm \
		--dns 10.96.0.10 --dns-search container-image-builder.svc.cluster.local \
		-e DOCKER_HOST=tcp://docker:2375 \
		-e CACHE_ENDPOINT_SERVER=minio \
		-e CACHE_ENDPOINT_PORT=9000 \
		-e CACHE_BUCKET=image-layers \
		-e CACHE_ACCESS_KEY=8Q9U4RBHKKB6HU70SRZ1 \
		-e CACHE_SECRET_KEY=oxxT2iqBlW6lgaDVe8ll6mP8z/OSVIUnn9cB4+Q0 \
		-v $(PWD):/project -w /project \
		utopiaplanitia/docker-image-builder-devtools:latest bats tests

.PHONY: queue-test
queue-test: ##@testing Runs more parallel builds then workers
	@docker build -q -f docker/devtools/Dockerfile -t utopiaplanitia/docker-image-builder-devtools:latest .
	./etc/queue-test.sh

.PHONY: deploy
deploy: ##@development Deploys the current code.
	docker build -f docker/devtools/Dockerfile   -t utopiaplanitia/docker-image-builder-devtools:latest .
	docker build -f docker/builder/Dockerfile    -t utopiaplanitia/docker-image-builder-worker:latest .
	docker build -f docker/dispatcher/Dockerfile -t utopiaplanitia/docker-image-builder-dispatcher:latest .
	kubectl apply -f kubernetes/namespace.yaml -f kubernetes
	./etc/restart-pods.sh