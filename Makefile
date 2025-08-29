build:
	CGO_ENABLED=0 go build -ldflags='-s -w' -o bin/quadlet-lsp

build_all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o "bin/quadlet-lsp-linux-amd64"
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags='-s -w' -o "bin/quadlet-lsp-linux-arm64"
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags='-s -w' -o "bin/quadlet-lsp-windows-arm64.exe"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o "bin/quadlet-lsp-windows-amd64.exe"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o "bin/quadlet-lsp-darwin-amd64"

test:
	go test -race ./...

