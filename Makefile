GO       ?= go
TOOL_MOD ?= -modfile tool.go.mod
TOOL     ?= $(GO) tool $(TOOL_MOD)

-include Makefile.local

test:
	$(GO) test -v -cover -failfast ./...

tools-install:
	$(GO) get $(TOOL_MOD) -tool github.com/mgechev/revive@latest
	$(GO) get $(TOOL_MOD) -tool github.com/segmentio/golines@latest
	$(GO) get $(TOOL_MOD) -tool mvdan.cc/gofumpt@latest
	$(GO) mod tidy $(TOOL_MOD)

format:
	find -iname '*.go' -print0 | xargs -0 $(TOOL) golines --max-len 80 -w -l
	find -iname '*.go' -print0 | xargs -0 $(TOOL) gofumpt -w -l

lint:
	$(TOOL) revive -config revive.toml -formatter friendly ./...

generate-enum:
	go run github.com/abice/go-enum@latest --marshal --no-iota \
		--output-suffix _generated -f ./dynamic/enums.go
