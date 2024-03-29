.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/details details/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/insert insert/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/search search/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
