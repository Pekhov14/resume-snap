# Makefile для Go проекта

# Переменные
GO=go
VENDOR_DIR=vendor

.PHONY: all build tidy vendor clean

all: build

build:
	$(GO) build -o resume cmd/main.go

# Dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod vendor

# Creacte vendor
vendor:
	$(GO) mod vendor

clean:
	rm -rf myapp $(VENDOR_DIR)

run:
	go run cmd/main.go

# todo
run-config:
	go run cmd/main.go --url=http://localhost:3000/cv --selector="#cv-container" --output="cv_anton_pekhov_backend_developer.pdf"