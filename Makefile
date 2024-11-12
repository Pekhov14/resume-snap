# Makefile для Go проекта

# Переменные
GO=go
VENDOR_DIR=vendor

# Цели
.PHONY: all build tidy vendor clean

# Основная цель
all: build

# Сборка проекта
build:
	$(GO) build -o resume .

# Удаление неиспользуемых зависимостей и создание папки vendor
tidy:
	$(GO) mod tidy
	$(GO) mod vendor

# Создание папки vendor
vendor:
	$(GO) mod vendor

# Очистка скомпилированных файлов и папки vendor
clean:
	rm -rf myapp $(VENDOR_DIR)

run-config:
	go run main.go

run-param:
	go run --url=https://anton-pekhov.vercel.app/cv --selector="#cv-container" --output="cv_anton_pekhov_backend_developer.pdf"