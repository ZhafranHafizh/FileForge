APP_NAME=fileforge

.PHONY: build run test fmt vet clean

build:
	go build -buildvcs=false -o bin/$(APP_NAME) .

run:
	go run .

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
