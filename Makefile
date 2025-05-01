fmt:
	go fmt ./...

vet: fmt
	go vet ./...

tidy: fmt
	go mod tidy

run/api: vet
	go run ./cmd/api
