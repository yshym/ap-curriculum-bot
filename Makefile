GOCMD=go
GOBUILD=GOOS=linux $(GOCMD) build -ldflags="-d -s -w" -tags netgo -installsuffix netgo

.PHONY: build clean test deploy logs

build:
	$(GOBUILD) -o bin/setwebhook setwebhook/main.go
	$(GOBUILD) -o bin/webhook webhook/main.go

clean:
	$(GOCMD) clean
	rm -rf ./bin

test:
	cd ./curriculum && $(GOCMD) test
	cd ./helpers && $(GOCMD) test

deploy: clean build test
	npm run deploy

logs:
	npm run logs
