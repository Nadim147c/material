GO ?= go

-include Makefile.local

test:
	$(GO) generate ./...
	$(GO) test -v -cover -failfast ./...

generate-enum:
	go run github.com/abice/go-enum@latest --marshal --noprefix \
		--output-suffix _generated -f ./dynamic/enums.go

	sed 's/X202/V202/' ./dynamic/enums_generated.go > ./dynamic/enums_generated.go.tmp && \
		mv ./dynamic/enums_generated.go.tmp ./dynamic/enums_generated.go
