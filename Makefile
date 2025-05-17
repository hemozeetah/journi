include .env
export

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

tidy: fmt
	go mod tidy

run/api: vet
	go run ./cmd/api

run/react:
	npm run --prefix cmd/api/v1/react/ dev

migrate: vet
	go run ./cmd/migrate

psql:
	@PGPASSWORD=${DATABASE.PASSWORD} \
		psql -U ${DATABASE.USER} -d ${DATABASE.NAME} -h ${DATABASE.HOST}
