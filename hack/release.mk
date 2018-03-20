
.PHONY: build-push
build-push: ##@release Build and push the images.
	docker build -f docker/devtools/Dockerfile   -t utopiaplanitia/docker-image-builder-devtools:latest .
	docker build -f docker/builder/Dockerfile    -t utopiaplanitia/docker-image-builder-worker:latest .
	docker build -f docker/dispatcher/Dockerfile -t utopiaplanitia/docker-image-builder-dispatcher:latest .
	docker push utopiaplanitia/docker-image-builder-devtools:latest
	docker push utopiaplanitia/docker-image-builder-worker:latest
	docker push utopiaplanitia/docker-image-builder-dispatcher:latest
