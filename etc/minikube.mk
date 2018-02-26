
.PHONY: minikube-start
minikube-start:
	$(MAKE) minikube-init
	$(MAKE) minikube-await

.PHONY: minikube-init
minikube-init:
	./etc/start-minikube.sh
	minikube update-context

.PHONY: minikube-await
minikube-await:
	./etc/await-minikube.sh

.PHONY: minikube-stop
minikube-stop:
	sudo -E systemctl stop localkube
	docker ps -aq --filter name=k8s | xargs -r docker rm -f

.PHONY: minikube-logs
minikube-logs:
	ktail -n container-image-builder

.PHONY: minikube-cli
minikube-cli:
	docker build -f docker/dev-tools/Dockerfile  -t utopiaplanitia/docker-image-builder-service:dev-tools-latest .
	docker run -ti --rm \
		--dns 10.96.0.10 --dns-search container-image-builder.svc.cluster.local \
		-e DOCKER_HOST=tcp://docker:2375 \
		-e CACHE_ENDPOINT_SERVER=minio \
		-e CACHE_ENDPOINT_PORT=9000 \
		-e CACHE_BUCKET=image-layers \
		-e CACHE_ACCESS_KEY=8Q9U4RBHKKB6HU70SRZ1 \
		-e CACHE_SECRET_KEY=oxxT2iqBlW6lgaDVe8ll6mP8z/OSVIUnn9cB4+Q0 \
		-v $(PWD):/project -w /project \
		utopiaplanitia/docker-image-builder-service:dev-tools-latest sh

.PHONY: minikube-tests
minikube-tests:
	docker build -f docker/dev-tools/Dockerfile  -t utopiaplanitia/docker-image-builder-service:dev-tools-latest .
	docker run -ti --rm \
		--dns 10.96.0.10 --dns-search container-image-builder.svc.cluster.local \
		-e DOCKER_HOST=tcp://docker:2375 \
		-e CACHE_ENDPOINT_SERVER=minio \
		-e CACHE_ENDPOINT_PORT=9000 \
		-e CACHE_BUCKET=image-layers \
		-e CACHE_ACCESS_KEY=8Q9U4RBHKKB6HU70SRZ1 \
		-e CACHE_SECRET_KEY=oxxT2iqBlW6lgaDVe8ll6mP8z/OSVIUnn9cB4+Q0 \
		-v $(PWD):/project -w /project \
		utopiaplanitia/docker-image-builder-service:dev-tools-latest bats tests

.PHONY: minikube-deploy
minikube-deploy: dispatcher builder
	docker build -f docker/builder/Dockerfile    -t utopiaplanitia/docker-image-builder-service:builder-latest .
	docker build -f docker/dispatcher/Dockerfile -t utopiaplanitia/docker-image-builder-service:dispatcher-latest .
	kubectl apply -f kubernetes/namespace.yaml -f kubernetes
	./etc/restart-pods.sh

.PHONY: minikube-deploy-multi-stage-build
minikube-deploy-multi-stage-build:
	docker build -f kubernetes/images/builder/Dockerfile    -t utopiaplanitia/docker-image-builder-service:builder-latest .
	docker build -f kubernetes/images/dispatcher/Dockerfile -t utopiaplanitia/docker-image-builder-service:dispatcher-latest .
	kubectl apply -f kubernetes/namespace.yaml -f kubernetes
	./etc/restart-pods.sh
