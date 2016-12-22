GO = go
pkgs = $(shell $(GO) list ./... | grep -v /vendor/)

.PHONY: all style format vet lint goimport test bench

all: style format vet lint goimport test bench

style:
	@echo ">> checking code style"
	@! gofmt -s -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

gocyclo:
	@echo ">> cyclomatic complexities"
	@gocyclo -over 15

lint:
	@echo ">> lint code"
	@golint -set_exit_status=1 ./...

goimport:
	@echo ">> goimports"
	@goimports -d ./

test:
	@echo ">> running tests"
	@SECURE_API_KEY=$(SECURE_API_KEY) \
	CONSUL_HTTP_ADDR=$(CONSUL_HTTP_ADDR) \
	$(GO) test -short -race $(pkgs)

bench:
	@echo ">> running benchmark"
	@SECURE_API_KEY=$(SECURE_API_KEY) \
	CONSUL_HTTP_ADDR=$(CONSUL_HTTP_ADDR) \
	$(GO) test -bench . -benchmem $(pkgs)
