SWAG ?= swag

.PHONY: docs main api

docs:
	$(SWAG) init -g api/main.go -o docs || \
		go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g api/main.go -o docs

api:
	go run api/main.go

dev: docs api
