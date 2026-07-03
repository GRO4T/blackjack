set quiet

PROTO_IN := "./blackjack.proto"
PROTO_OUT_DIR := "./bjack-api/proto"
PROTO_OUT := "./bjack-api/proto/blackjack.pb.go ./bjack-api/proto/blackjack_grpc.pb.go"
API_DIR := "./bjack-api"
API_EXECUTABLE := "./bin/bjack-api"
UI_DIR := "./bjack-ui"
CI_IMAGE_TAG := "0.0.5"

default:
  just --list

[no-quiet]
[group("setup")]
setup: setup_api setup_ui

[group("setup")]
setup_api:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

[group("setup")]
setup_ui:
	cd {{UI_DIR}} && npm install

[no-quiet]
clean:
	-rm {{PROTO_OUT}}
	-rm {{API_DIR}}/{{API_EXECUTABLE}}
	-rm -r {{UI_DIR}}/dist

# compile gRPC
[group("proto")]
proto:
	protoc -I=. --go_out={{PROTO_OUT_DIR}} --go_opt=paths=source_relative \
		--go-grpc_out={{PROTO_OUT_DIR}} --go-grpc_opt=paths=source_relative \
		{{PROTO_IN}}

# build API (native)
[group("api")]
build_api: proto
	cd {{API_DIR}} && go build -o {{API_EXECUTABLE}}

# run API locally
[group("api")]
run_api_dev: build_api
	eval $(cat {{API_DIR}}/.env.development) {{API_DIR}}/{{API_EXECUTABLE}}

# build API (docker) 
[group("api")]
build_api_image:
	docker build -t bjack-api -f docker/api.Dockerfile .

# run API locally (docker)
[group("api")]
run_api_image_dev: build_api_image
	docker run -p 8000:8000 --env-file {{API_DIR}}/.env.development bjack-api

# npm build
[group("ui")]
build_ui:
	npm run build

# start UI server (MODE="dev|preview")
[group("ui"), arg("MODE", pattern="dev|preview")]
run_ui MODE: build_ui
	#!/usr/bin/env bash
	cd {{UI_DIR}}
	if [[ "{{MODE}}" == "dev" ]]; then
		npm run dev
	elif [[ "{{MODE}}" == "preview" ]]; then
		npm run preview
	else
		echo "Unknown mode: {{MODE}}"
	fi

[group("validation")]
test: proto
	cd {{API_DIR}} && go test ./...

[no-quiet]
[group("validation")]
lint: proto
	cd {{API_DIR}} && golangci-lint run
	cd {{UI_DIR}} && npm run lint

[no-quiet]
[group("validation")]
fmt: proto
	cd {{API_DIR}} && test -z $(gofmt -l .) || gofmt -l . | false
	cd {{UI_DIR}} && npx prettier . --check

[no-quiet]
[group("validation")]
fmt_fix: proto
	cd {{API_DIR}} && gofmt -s -w .
	cd {{UI_DIR}} && npx prettier . --write

# build docker image used in GitHub Workflows
[group("ci_image")]
build_ci_image:
	docker build -t dkolaska/blackjack-ci:{{CI_IMAGE_TAG}} -f docker/ci.Dockerfile --platform linux/amd64 .

# push docker image used in GitHub Workflows
[group("ci_image")]
push_ci_image:
	docker push dkolaska/blackjack-ci:{{CI_IMAGE_TAG}}
