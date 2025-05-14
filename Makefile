GO ?= go

-include Makefile.local

TEST := $(GO) test -v -failfast
test:
	$(TEST) ./color
	$(TEST) ./num
	$(TEST) ./quantizer
