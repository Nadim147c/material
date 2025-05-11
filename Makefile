GO ?= go

include Makefile.local

build:
	@echo "build doesn't exists"

test:
	$(GO) test -v color num
