
.PHONY: up
up: builder dispatcher
	docker-compose up --build -d --remove-orphans
	$(MAKE) logs

.PHONY: down
down:
	docker-compose down --remove-orphans -t 1

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: cli
cli:
	docker-compose exec dev-tools sh

.PHONY: tests
tests:
	docker-compose exec dev-tools bats tests
