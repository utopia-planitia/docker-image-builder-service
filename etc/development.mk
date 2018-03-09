
.PHONY: cli
cli: .docker-image-devtools ##@development Opens a command line interface with development tools.
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

.PHONY: deploy
deploy: .docker-image-devtools .docker-image-builder .docker-image-dispatcher ##@development Deploys the current code.
	kubectl apply -f kubernetes/namespace.yaml -f kubernetes
	./etc/restart-pods.sh

.docker-image-devtools: caching-scripts/entrypoint.sh docker/devtools/Dockerfile
	docker build -f docker/devtools/Dockerfile  -t utopiaplanitia/docker-image-builder-devtools:latest .
	touch .docker-image-devtools

.docker-image-builder: $(shell find caching-scripts cmd/builder -type f) docker/builder/Dockerfile
	docker build -f docker/builder/Dockerfile  -t utopiaplanitia/docker-image-builder-worker:latest .
	touch .docker-image-builder

.docker-image-dispatcher: $(shell find cmd/dispatcher -type f) docker/dispatcher/Dockerfile
	docker build -f docker/dispatcher/Dockerfile  -t utopiaplanitia/docker-image-builder-dispatcher:latest .
	touch .docker-image-dispatcher
