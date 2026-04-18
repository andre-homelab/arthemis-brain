SWAG ?= swag

.PHONY: docs main

docs:
	$(SWAG) init -g inputs/api/main.go -o docs || \
		go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g inputs/api/main.go -o docs

main:
	go run cmd/api/main.go

dev: docs api 
