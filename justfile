set quiet

PROTOC_VERSION := "35.1"
PROTO_IN := "./blackjack.proto"
PROTO_OUT_DIR := "./bjack-api/proto"
PROTO_OUT := "./bjack-api/proto/blackjack.pb.go ./bjack-api/proto/blackjack_grpc.pb.go"
API_DOCKER_IMAGE := "golang:1.26.4-alpine3.24"
API_DIR := "./bjack-api"
API_EXECUTABLE := "./bin/bjack-api"
UI_DIR := "./bjack-ui"
CI_IMAGE_TAG := "0.0.6"

default:
  just --list

[no-quiet]
[group("setup")]
setup: setup_api setup_ui

[group("setup")]
setup_api:
    just _install_grpc
    just _install_golangci_lint

_install_grpc:
	#!/usr/bin/env sh
	set -eu
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2

	OS_NAME="linux"
	ARCH="$(uname -m)"

	case "${ARCH}" in
	  x86_64|amd64)
	    ARCH_NAME="x86_64"
	    ;;
	  arm64|aarch64)
	    ARCH_NAME="aarch_64"
	    ;;
	  *)
	    echo "Unsupported architecture: ${ARCH}" >&2
	    exit 1
	    ;;
	esac

	ZIP_FILE="protoc-{{PROTOC_VERSION}}-${OS_NAME}-${ARCH_NAME}.zip"

	URL="https://github.com/protocolbuffers/protobuf/releases/download/v{{PROTOC_VERSION}}/${ZIP_FILE}"
	echo "Downloading protoc from ${URL}..."

	# Download using wget or curl
	TMP_ZIP=$(mktemp)
	if command -v wget >/dev/null 2>&1; then
	  wget -O "${TMP_ZIP}" "${URL}"
	elif command -v curl >/dev/null 2>&1; then
	  curl -L -o "${TMP_ZIP}" "${URL}"
	else
	  echo "Error: neither wget nor curl is installed." >&2
	  rm -f "${TMP_ZIP}"
	  exit 1
	fi

	echo "Installing protoc to ${HOME}/.local..."
	unzip -o "${TMP_ZIP}" -d "${HOME}/.local"
	rm -f "${TMP_ZIP}"
	echo "protoc installed successfully."

_install_golangci_lint:
    wget -O- -nv https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2

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
	docker build --build-arg "VERSION={{API_DOCKER_IMAGE}}" -t bjack-api -f docker/api.Dockerfile .

# run API locally (docker)
[group("api")]
run_api_image_dev: build_api_image
	docker run -p 8000:8000 --env-file {{API_DIR}}/.env.development bjack-api

# npm build
[group("ui")]
build_ui:
	cd {{UI_DIR}} && npm run build


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
	docker build --build-arg "VERSION={{API_DOCKER_IMAGE}}" -t dkolaska/blackjack-ci:{{CI_IMAGE_TAG}} -f docker/ci.Dockerfile --platform linux/amd64 .

# push docker image used in GitHub Workflows
[group("ci_image")]
push_ci_image:
	docker push dkolaska/blackjack-ci:{{CI_IMAGE_TAG}}
