.PHONY: up down build logs logs-db ps status tables restart sqlc help generate

include .env
export

# Variables
SWAGGER_FILE=api/swagger.json
GENERATED_DIR=internal/api/generated
GENERATED_FILE=$(GENERATED_DIR)/api.gen.go

up:
	docker-compose up -d

build:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

logs-db:
	docker-compose logs -f ari_pari_db

ps:
	docker-compose ps

status:
	docker-compose exec postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB)

tables:
	docker-compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "\dt"

restart: down up

sqlc:
	sqlc generate

help:
	@echo "Доступные команды:"
	@echo "  make up        - Старт всех сервисов"
	@echo "  make build     - Пересобрать и запустить"
	@echo "  make down      - Остановка сервисов"
	@echo "  make logs      - Показать логи"
	@echo "  make logs-db   - Показать логи PostgreSQL"
	@echo "  make ps        - Показать запущенные контейнеры"
	@echo "  make status    - Проверить статус БД"
	@echo "  make tables    - Показать таблицы в БД"
	@echo "  make restart   - Перезапустить сервисы"
	@echo "  make sqlc   	- Сгенерировать go код для CRUD запросов"

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
			-generate types,chi-server \
			-package api \
			-o $(GENERATED_FILE) \
			$(SWAGGER_FILE)

mockgen:
	mockgen -source=$(source) -package=mocks -destination=$(destination)
