genmocks:
	go install github.com/golang/mock/mockgen@latest
	go generate ./...

tests: genmocks
	go test ./...

build-server-for-windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o ./build/server-windows.exe ./cmd/client

build-client-for-windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o ./build/client-windows.exe ./cmd/client

build-server-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/server-macos ./cmd/server

build-client-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/client-macos ./cmd/client

build-both-for-windows: build-server-for-windows build-client-for-windows

build-both-for-macos: build-server-for-macos build-client-for-macos

