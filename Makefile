GO ?= go

include Makefile.local

test:
	$(GO) test -v ./color ./num ./quantizer
