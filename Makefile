# Makefile для создания миграций

# Переменные которые будут использоваться в наших командах (Таргетах)
DB_DSN := "postgres://postgres:yourpassword@localhost:5432/main?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Таргет для создания новой миграции
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# Применение миграций
migrate:
	$(MIGRATE) up

# Откат миграций
migrate-down:
	$(MIGRATE) down

# для удобства добавим команду run, которая будет запускать наше приложение
run:
	go run cmd/app/main.go # Теперь при вызове make run мы запустим наш сервер

# это цель (target) в Makefile, которая используется для генерации кода на основе OpenAPI-спецификации с помощью
# oapi-codegen
gen:
	oapi-codegen -config openapi/.openapi -include-tags tasks -package tasks openapi/openapi.yaml \
	> ./internal/web/tasks/api.gen.go

# Линтер - инструмент, который анализирует нашу кодовую базу и указывает на все места в которых есть ошибки
lint:
	golangci-lint run --out-format=colored-line-number