GO     ?= go
REVIVE ?= revive

-include Makefile.local

test:
	$(GO) test -v -cover -failfast ./...

lint:
	$(REVIVE) -config revive.toml -formatter friendly ./...

generate-enum:
	go run github.com/abice/go-enum@latest --marshal --no-iota \
		--output-suffix _generated -f ./dynamic/enums.go
