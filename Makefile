genmocks:
	go install github.com/golang/mock/mockgen@latest
	go generate ./...

tests: genmocks
	go test ./...

build-server-for-windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o ./build/server.exe ./cmd/server

build-client-for-windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o ./build/client.exe ./cmd/client

build-both-for-windows: build-server-for-windows build-client-for-windows

