include .env

BIN_FILENAME := banner_rotator
BIN := "./bin/$(BIN_FILENAME)"
DOCKER_COMPOSE_DEV_PATH="./deployments/docker-compose.dev.yaml"
DOCKER_COMPOSE_PATH="./deployments/docker-compose.yaml"
DOCKER_COMPOSE_TEST_PATH="./deployments/docker-compose.test.yaml"

build-binary:
	go build -v -o $(BIN) .

build:
	docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_DEV_PATH) build
	go build -v -o $(BIN) .

up:
	docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) up -d
	while ! docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) exec -T --user postgres db psql -c "select 'db ready!'" > /dev/null; do sleep 1; done;
	while ! curl -f -s http://localhost:15672 > /dev/null; do sleep 1; done;

down:
	docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) down --remove-orphans

install-migrator:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: install-migrator
	goose --dir ./migrations postgres postgres://postgres:password@localhost:${POSTGRES_PORT}/postgres?sslmode=disable up

run:
	docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_DEV_PATH) up -d

stop:
	docker-compose --env-file .env -p "$(BIN_FILENAME)_dev" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_DEV_PATH) down --remove-orphans

test:
	go test -race -v -count 100 ./...

install-lint-deps:
ifeq (,$(wildcard ./bin/golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.42.0
endif

lint: install-lint-deps
	./bin/golangci-lint run --config .golangci.yaml --color always ./...

install-protoc:
	apt install -y protobuf-compiler

generate: install-protoc
	protoc \
	--proto_path=./api/ --go_out=./internal/server/pb --go-grpc_out=./internal/server/pb \
	api/*.proto
	go generate ./...

integration-tests: install-migrator
	docker-compose --env-file .env.test -p "$(BIN_FILENAME)_test" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_TEST_PATH) up --build -d
	while ! docker-compose --env-file .env.test -p "$(BIN_FILENAME)_test" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_TEST_PATH) exec -T --user postgres db psql -c "select 'db ready!'" > /dev/null; do sleep 1; done;
	while ! curl -f -s http://localhost:15677 > /dev/null; do sleep 1; done;
	goose --dir ./migrations postgres postgres://postgres:password@localhost:5437/postgres?sslmode=disable up
	go test -v ./... --tags integration; \
	e=$$?; \
	docker-compose --env-file .env.test -p "$(BIN_FILENAME)_test" -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_TEST_PATH) down --remove-orphans; \
	exit $$e

.PHONY: build up down migrate run stop test lint generate integration-tests