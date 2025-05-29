test:
	go test ./...

lint:
	go vet ./... && go fmt -s ./...

run:
	go run ./cmd/main.go
