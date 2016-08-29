install:
	@echo "Not implemented"

build:
	@echo "Not implemented"

deps:
	go get -u github.com/nats-io/nats
	go get -u github.com/smartystreets/goconvey/convey
	go get -u github.com/golang/lint/golint
	go get -u github.com/jwilder/dockerize

test:
	go test -v

lint:
	golint ./...
	go vet ./...
