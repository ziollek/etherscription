PATH := $(PATH):/go/bin:$(PWD)/bin/dependencies
BINARY_NAME=etherscription
BINARY_FILE := ./bin/$(BINARY_NAME)
GOTEST_FLAGS := -cover -race -v -count=1 -timeout 60s
NODE="http://localhost:8545"
PORT=8888


.PHONY: test run
build:
	@echo ">> building application"
	go build -trimpath \
	-o $(BINARY_FILE) \
	./cmd/etherscription

run:
	go run ./cmd/etherscription --node $(NODE) --port $(PORT)

race-run:
	go run -race ./cmd/etherscription --node $(NODE) --port $(PORT)
	
test:
	@echo ">> running all tests"
	go test -short $(GOTEST_FLAGS) ./...

gofumpt:
	gofumpt --extra -w .

clean:
	go clean
	rm $(BINARY_FILE)
