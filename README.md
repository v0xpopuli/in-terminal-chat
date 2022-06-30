# in-terminal-chat
[![.github/workflows/tests.yml](https://github.com/v0xpopuli/in-terminal-chat/actions/workflows/tests.yml/badge.svg?branch=development)](https://github.com/v0xpopuli/in-terminal-chat/actions/workflows/tests.yml)
![Coverage](https://img.shields.io/badge/Coverage-90.9%25-brightgreen)

<p>Simple chat that works in terminal. <br/>
Just small project for keeping skills sharp and learn something new. <br/><p>

Run **server** -> `go run cmd/server/main.go` \
Run **client** -> `go run cmd/client/main.go`

`Usage of server:` \
&nbsp;&nbsp;&nbsp;&nbsp;`-address string`
`http service address (default "localhost:8080")` \
&nbsp;&nbsp;&nbsp;&nbsp;`-debug`
`debug mode (default true)`

`Usage of client:` \
&nbsp;&nbsp;&nbsp;&nbsp;`-address string`
`http service address (default "ws://localhost:8080")` \
&nbsp;&nbsp;&nbsp;&nbsp;`-name string`
`client name`