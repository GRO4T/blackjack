set quiet

PROTO_IN := "./blackjack.proto"
PROTO_OUT_DIR := "./bjack-api/proto"
PROTO_OUT := "./bjack-api/proto/blackjack.pb.go ./bjack-api/proto/blackjack_grpc.pb.go"
API_DIR := "./bjack-api"
API_EXECUTABLE := "./bin/bjack-api"
UI_DIR := "./bjack-ui"
CI_IMAGE_TAG := "0.0.4"

default:
  just --list

setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	cd bjack-ui && npm install

clean:
	rm -f {{PROTO_OUT}} {{API_DIR}}/{{API_EXECUTABLE}}

proto:
	protoc -I=. --go_out={{PROTO_OUT_DIR}} --go_opt=paths=source_relative \
		--go-grpc_out={{PROTO_OUT_DIR}} --go-grpc_opt=paths=source_relative \
		{{PROTO_IN}}

[group("api")]
build_api: proto
	cd {{API_DIR}} && go build -o {{API_EXECUTABLE}}

[group("api")]
run_api: build_api
	{{API_DIR}}/{{API_EXECUTABLE}}

[group("ui")]
run_ui MODE:
	#!/usr/bin/env bash
	cd {{UI_DIR}}
	if [[ "{{MODE}}" == "dev" ]]; then
		npm run dev
	else
		echo "Unknown mode: {{MODE}}"
	fi

[group("validation")]
test: proto
	cd {{API_DIR}} && go test ./...

[group("validation")]
lint: proto
	cd {{API_DIR}} && golangci-lint run

[group("validation")]
fmt: proto
	cd {{API_DIR}} && test -z $(gofmt -l .) || gofmt -l . | false

[group("validation")]
fmt_fix: proto
	cd {{API_DIR}} && gofmt -s -w .

[group("ci_image")]
build_ci_image:
	docker build -t dkolaska/blackjack-ci:{{CI_IMAGE_TAG}} -f ci.dockerfile --platform linux/amd64 .

[group("ci_image")]
push_ci_image:
	docker push dkolaska/blackjack-ci:{{CI_IMAGE_TAG}}
