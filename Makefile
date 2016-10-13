install:
	@echo "Not implemented"

build:
	@echo "Not implemented"

deps:
	go get github.com/nats-io/nats
	go get github.com/smartystreets/goconvey/convey
	go get github.com/golang/lint/golint
	go get github.com/jwilder/dockerize

test:
	go test -v

lint:
	golint ./...
	go vet ./...
