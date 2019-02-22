get:
	@echo ">> Getting any missing dependencies.."
	go get -t ./...
.PHONY: get

install:
	go install github.com/jstolp/pofadder-go
.PHONY: install

run: install
	./pofadder-go server
.PHONY: run

test:
	go test ./...
.PHONY: test

do:
	@echo ">> Doing it for ya, master..."
	go build -v
	@echo ">> build completed"
	./pofadder-go server
.PHONY: do

fmt:
	@echo ">> Running Gofmt.."
	gofmt -l -s -w .
