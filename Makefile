APP       := api
CMD_DIR   := cmd/$(APP)/main.go
BIN_DIR   := bin
BIN       := $(BIN_DIR)/$(APP)
PKG       := ./...
TEST_PKGS := ./tests/...

.DEFAULT_GOAL := help
.PHONY: help build run test fmt vet tidy clean

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build   Build the api binary"
	@echo "  run     Build and run the binary"
	@echo "  test    Run tests"
	@echo "  fmt     Format source code"
	@echo "  vet     Run go vet"
	@echo "  tidy    Run go mod tidy"
	@echo "  clean   Remove build artifacts"

build:
	@go build -o $(BIN) $(CMD_DIR)

run: build
	@./$(BIN)

test:
	@go test -v $(TEST_PKGS)

fmt:
	@go fmt $(PKG)

vet:
	@go vet $(PKG)

tidy:
	@go mod tidy

clean:
	@rm -rf $(BIN_DIR)
