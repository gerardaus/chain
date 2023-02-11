unit-test:
	go test -tags unit ./...

test:
	go test -tags all ./...

integration-test:
	go test -tags integration ./...

build-examples:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/scarlett examples/scarlett/main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/ai examples/ai/main.go
