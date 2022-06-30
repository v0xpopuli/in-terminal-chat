genmocks:
	go install github.com/golang/mock/mockgen@latest
	go generate ./...

tests: genmocks
	go test ./...

build-server-for-windows:
	GOOS=windows GOARCH=386 go build -o ./build/windows/server.exe ./cmd/server

build-client-for-windows:
	GOOS=windows GOARCH=386 go build -o ./build/windows/client.exe ./cmd/client

build-server-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/macos/server ./cmd/server

build-client-for-macos:
	GOOS=darwin GOARCH=amd64 go build -o ./build/macos/client ./cmd/client

build-both-for-windows: build-server-for-windows build-client-for-windows

build-both-for-macos: build-server-for-macos build-client-for-macos

