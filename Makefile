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
	$(GO) build -o resume cmd/main.go

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

run:
	go run cmd/main.go --url=http://localhost:3000/cv --selector="#cv-container" --output="cv_anton_pekhov_backend_developer.pdf" --HeightAdjustment=-300 --scale=0.8


run-config:
	go run cmd/main.go --url=http://localhost:3000/cv --selector="#cv-container" --output="cv_anton_pekhov_backend_developer.pdf"