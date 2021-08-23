test: tidy
	go test -v -race ./...

cover: tidy
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

tidy:
	go mod tidy

fmt:
	goimports -l -w .
