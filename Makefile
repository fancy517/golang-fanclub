.PHONY: build clean tool lint help

DB_DSN="mysql://root:root@tcp(localhost:3306)/fanclub"
VERSION=0.1.0
BINARY_NAME=sob

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

all: build

build:
	@go build -o ./bin/${BINARY_NAME} -v ./cmd

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golangci-lint run

clean:
	rm -rf ./bin

dev:
	@go run ./cmd

prod:
	mode=prod go run ./cmd

# *****************************************************************
# **************************** Database ***************************
# *****************************************************************

migration_new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

migration_up: confirm
	migrate -path ./migrations -database ${DB_DSN} -verbose up

migration_down: confirm
	migrate -path ./migrations -database ${DB_DSN} -verbose down

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'