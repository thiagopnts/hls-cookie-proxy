BINARY ?= hls-cookie-proxy

$(BINARY): *.go
	@go build -o $(BINARY)

.PHONY: run
run: $(BINARY)
	@./$(BINARY)

.PHONY: clean
clean:
	@rm -rf $(BINARY)

.PHONY: all
all: $(BINARY)

.PHONY: test
test:
	go test
