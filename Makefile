build:
	CGO_ENABLED=0 go build -o bin/quadlet-lsp

build_all:
	goreleaser release --clean --skip=publish --skip=validate

test:
	go test -race ./...

