SHELL := /bin/bash
TOOLS_BIN := tools/bin
PILTOVER := $(TOOLS_BIN)/piltover

.PHONY: tools
tools: $(PILTOVER)

$(PILTOVER):
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
	@cd tools && go test ./...

.PHONY: clean
clean:
	@rm -rf $(TOOLS_BIN)
