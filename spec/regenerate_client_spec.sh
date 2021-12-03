# Download the latest mainline rpc-spec.
wget https://raw.githubusercontent.com/pokt-network/pocket-core/staging/doc/specs/rpc-spec.yaml -O spec/rpc-spec.yaml

# Generate go client.
oapi-codegen --config=spec/oapi_codegen.yaml spec/rpc-spec.yaml > spec/pocket_client.gen.go

# TODO: Consider using openapi-generator or something else (https://gist.github.com/craigmurray1120/8e87d88a076d49ec9c43636a313cfa66) instead.
# Configurations options: https://github.com/OpenAPITools/openapi-generator/blob/master/docs/generators/go.md
# openapi-generator generate -g go -i rpc-spec.yaml -o specv --package-name specv