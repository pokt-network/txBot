module github.com/pokt-network/txbot

go 1.13

replace github.com/tendermint/tendermint => github.com/pokt-network/tendermint v0.32.11-0.20210427155510-04e1c67f3eed // indirect

require (
	github.com/deepmap/oapi-codegen v1.9.0
	github.com/pokt-network/pocket-core v0.0.0-20210429190449-f794bc74b167
	github.com/tjarratt/babble v0.0.0-20210505082055-cbca2a4833c1
)