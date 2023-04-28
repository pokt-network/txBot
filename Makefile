.SILENT:

.PHONY: help
## This help screen
help:
	printf "Available targets\n\n"
	awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "%-30s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: swagger_check
swagger_check: # Check if swagger is installed
	{ \
	if ( ! ( command -v swagger >/dev/null ) ); then \
		echo "Seems like you don't have swagger installed.\nRun 'brew install go-swagger' on macOS.\nAlternatively, see https://goswagger.io/install.html or https://swagger.io/docs/open-source-tools/swagger-codegen/ for more details"; \
		exit 1; \
	fi; \
	}


.PHONY: start
## Compile & start a tx-bot instance
start:
	go build && ./txbot

.PHONY: start_burst
## Compile & start a tx-bot instance in burset most
start_burst:
	go build
	TX_CONFIG_FILE="burst_config.json" ./txbot

.PHONY: gen_client_spec
## Regenerate the go client based on the rpc spec
gen_client_spec: swagger_check
	wget https://raw.githubusercontent.com/pokt-network/pocket-core/staging/doc/specs/rpc-spec.yaml -O spec/rpc-spec.yaml
	oapi-codegen --config=spec/oapi_codegen.yaml spec/rpc-spec.yaml > spec/pocket.gen.raw.go