
.PHONY: up
up:
	docker-compose up --build -d

.PHONY: down
down:
	docker-compose down --remove-orphans

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: cli
cli:
	docker-compose exec dev-tools sh

.PHONY: tests
tests:
	docker-compose exec dev-tools bats tests
