.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
	go test -v -race -timeout 30s ./..

.PHONY: migrateTable
	migrate -path migrations -database "postgres://localhost/restapi_dev?port=5432&user=postgres&password=7777&sslmode=disable" up

.DEFAULT_GOAL := build