GO ?= go

-include Makefile.local

test:
	$(GO) test -v -failfast ./color ./num ./quantizer
