BIN := bin/cloudcanal

.PHONY: build test clean

build:
	mkdir -p $(dir $(BIN))
	go build -o $(BIN) ./cmd/cloudcanal

test:
	go test ./...

clean:
	rm -rf $(dir $(BIN))
