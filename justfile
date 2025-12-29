test:
    gotestsum --format pkgname-and-test-fails

bench:
    go test -bench=. -run=^$ ./...

lint:
    revive -config revive.toml ./...

generate:
    go-enum --marshal --no-iota --output-suffix _generated -f ./dynamic/enums.go
    go run ./scripts/enum_alias.go ./dynamic
    gofumpt -w ./enum_generated_alias.go
