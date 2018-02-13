
.PHONY: up
up:
	docker-compose up --build --scale worker=2

.PHONY: down
down:
	docker-compose down --remove-orphans
