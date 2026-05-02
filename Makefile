SHELL := /bin/bash
TOOLS_BIN := tools/bin
PILTOVER := $(TOOLS_BIN)/piltover

.PHONY: tools
tools:
	@mkdir -p $(TOOLS_BIN)
	@echo "→ [.] $$ go build -o $(PILTOVER) ./tools/cmd/piltover"
	@cd tools && go build -o ../$(PILTOVER) ./cmd/piltover

.PHONY: doctor
doctor: tools
	@$(PILTOVER) doctor

.PHONY: ls
ls: tools
	@$(PILTOVER) ls

.PHONY: ci
ci: tools
	@$(PILTOVER) ci

.PHONY: test
test:
	@cd tools && go test -race -count=1 ./...

.PHONY: lint
lint:
	@cd tools && golangci-lint run ./...

.PHONY: vet
vet:
	@cd tools && go vet ./...

.PHONY: verify
verify: lint vet test

.PHONY: clean
clean:
	@rm -rf $(TOOLS_BIN)
