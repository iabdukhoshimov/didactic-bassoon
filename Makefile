CURRENT_DIR=$(shell pwd)

-include .env

PSQL_CONTAINER_NAME?=postgres-container
PROJECT=$(shell basename ${CURRENT_DIR})
PSQL_URI?=postgres://postgres:postgres@localhost:5432/go_auth_service?sslmode=disable
APP_CMD_DIR=${CURRENT_DIR}/cmd
NAME=alter_organization_key_nullable

TAG=latest


.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: createdb
createdb:
	docker exec -it ${PSQL_CONTAINER_NAME} createdb -U postgres ${PROJECT_NAME}

.PHONY: execdb
execdb:
	docker exec -it ${PSQL_CONTAINER_NAME} psql -U postgres ${PROJECT_NAME}

.PHONY: dropdb
dropdb:
	docker exec -it ${PSQL_CONTAINER_NAME} dropdb -U postgres ${PROJECT_NAME}

.PHONY: execdb
cleandb:
	docker exec -it ${PSQL_CONTAINER_NAME} psql -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" ${PSQL_URI}

.PHONY: migrate_up
migrate_up:
	migrate -path ./migrations/postgres -database "${PSQL_URI}" -verbose up

.PHONY: migrate_down
migrate_down:
	migrate -path ./migrations/postgres -database "${PSQL_URI}&x-migrations-table=migrations_${PROJECT}" -verbose down

.PHONY: migrate_status
migrate_status:
	migrate -path ./migrations/postgres -database "${PSQL_URI}&x-migrations-table=migrations_${PROJECT}" -verbose status

.PHONY: migrate_create
migrate_create:
	migrate create -ext sql -dir ./migrations/postgres -seq ${NAME}

.PHONY: pull_submodules
pull_submodules:
	git submodule update --init --recursive

.PHONY: update_submodules
update_submodules:
	git submodule update --recursive --remote

build_image:
	docker build --rm -t "${REGISTRY_NAME}/${PROJECT_NAME}:${TAG}" .

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${PROJECT} ${APP_CMD_DIR}/main.go

.PHONY: dev_environment_start
dev_environment_start:
	docker compose -f docker-compose.yml up -d --force-recreate

.PHONY: dev_environment_stop
dev_environment_stop:
	docker compose -f docker-compose.yml down

.PHONY: dev_environment_remove
dev_environment_remove:e
	docker compose -f docker-compose.yml down --volumes

push_image:
	docker push "${REGISTRY_NAME}/${PROJECT}:${TAG}"

proto:
	rm -f generated/**/*.go
	rm -f doc/swagger/*.swagger.json
	mkdir -p generated
	protoc \
		--proto_path=tr_protos --go_out=generated --go_opt=paths=source_relative \
		--go-grpc_out=generated --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=generated --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger_docs,use_allof_for_refs=true,disable_service_tags=false \
			tr_protos/**/*.proto
