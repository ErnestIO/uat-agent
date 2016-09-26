install:
	@echo "Not implemented"

build:
	@echo "Not implemented"

deps:
	go get github.com/nats-io/nats
	go get github.com/smartystreets/goconvey/convey
	go get github.com/golang/lint/golint
	go get github.com/jwilder/dockerize
	go get github.com/gucumber/gucumber/cmd/gucumber

test:
	go test -v
	gucumber

lint:
	golint ./...
	go vet ./...
