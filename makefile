GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

build:
	GO111MODULE=on $(GOBUILD) -v ./cmd/cli/moonapi.go

test:
	GO111MODULE=on $(GOTEST) -v ./...

clean:
	rm *.exe *.zip main