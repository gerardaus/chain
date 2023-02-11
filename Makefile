unit-test:
	go test -tags unit ./...

test:
	go test -tags all ./...

integration-test:
	go test -tags integration ./...

build-examples:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o examples/scarlett/main examples/scarlett/main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o examples/ai/main examples/ai/main.go
