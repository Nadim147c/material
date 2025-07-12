GO ?= go

-include Makefile.local

test:
	$(GO) test -v -cover -failfast ./...
