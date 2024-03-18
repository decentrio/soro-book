set -eo pipefail

buf generate --path="./proto/api" --template="buf.gen.yaml" --config="buf.yaml"