OSNAME=$(shell go env GOOS)
TMP_DIR = /tmp
GOLANGCI_LINT_VERSION=1.46.2
GOLANGCI_DIR = $(TMP_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_TMP_BIN = $(GOLANGCI_DIR)/golangci-lint
GOLANGCI_LINT_ARCHIVE = golangci-lint-$(GOLANGCI_LINT_VERSION)-$(OSNAME)-amd64.tar.gz
BUILD_OUTPUT_DIR= artifacts
BINARY_NAME= $(BUILD_OUTPUT_DIR)/svc

.PHONY: prep
prep: deps lint unit_test

.PHONY: deps
deps:
	go mod tidy
	go mod verify

.PHONY: lint
lint: $(GOLANGCI_TMP_BIN)
	$(GOLANGCI_DIR)/golangci-lint run --allow-parallel-runners --tests=false --modules-download-mode=mod

# install a local golangci-lint if not found.
$(GOLANGCI_TMP_BIN):
	curl -OL https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/$(GOLANGCI_LINT_ARCHIVE)
	mkdir -p $(GOLANGCI_DIR)/
	tar -xf $(GOLANGCI_LINT_ARCHIVE) --strip-components=1 -C $(GOLANGCI_DIR)/
	chmod +x $(GOLANGCI_TMP_BIN)
	rm -f $(GOLANGCI_LINT_ARCHIVE)

.PHONY: unit_test
unit_test:
	go test -v -cover -coverprofile=coverage.out -count=1 ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BUILD_OUTPUT_DIR)/*.so

.PHONY: build
build: clean
	CGO_ENABLED=1 \
	GOOS=$(OSNAME) \
	go build \
	-o $(BINARY_NAME) \
	.
