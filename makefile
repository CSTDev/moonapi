GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

build:
	GO111MODULE=on $(GOBUILD) -v ./cmd/cli/moonapi.go

build-downloader:
	GO111MODULE=on $(GOBUILD) -v ./cmd/downloader/downloader.go

test:
	GO111MODULE=on $(GOTEST) -v ./...

clean:
	rm *.exe *.zip main