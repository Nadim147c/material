GO ?= go

-include Makefile.local

TEST := $(GO) test -v -failfast
test:
	# $(TEST) ./.
	$(TEST) ./color
	# $(TEST) ./num
	# $(TEST) ./score
	# $(TEST) ./blend
	# $(TEST) ./contrast
	# $(TEST) ./dislike
	# $(TEST) ./temperature
	# $(TEST) ./dynamic
	# $(TEST) ./quantizer
