module github.com/opsorch/opsorch-slack-adapter

go 1.22

require (
	github.com/opsorch/opsorch-core v0.0.2
	github.com/slack-go/slack v0.12.5
)

require github.com/gorilla/websocket v1.4.2 // indirect

replace github.com/opsorch/opsorch-core => ../opsorch-core
