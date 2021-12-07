# Download the latest mainline rpc-spec.
wget https://raw.githubusercontent.com/pokt-network/pocket-core/staging/doc/specs/rpc-spec.yaml -O spec/rpc-spec.yaml

# Generate go client.
oapi-codegen --config=spec/oapi_codegen.yaml spec/rpc-spec.yaml > spec/pocket_client.gen.go