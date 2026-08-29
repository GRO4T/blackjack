PROTOC_VERSION := "35.1"
PROTO_IN := "./blackjack.proto"
PROTO_OUT_DIR := "./proto"
PROTO_OUT := "./proto/blackjack.pb.go ./proto/blackjack_grpc.pb.go"
API_DOCKER_IMAGE := "golang:1.26.4-alpine3.24"
API_EXECUTABLE := "./bin/bjack-api"
DOCKER_CMD := if env_var_or_default("CONTAINER_ID", "") != "" { "distrobox-host-exec podman" } else { "docker" }

default:
    just --list api

setup: _install_grpc _install_golangci_lint

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

# compile gRPC
proto:
	protoc -I=. --go_out={{PROTO_OUT_DIR}} --go_opt=paths=source_relative \
		--go-grpc_out={{PROTO_OUT_DIR}} --go-grpc_opt=paths=source_relative \
		{{PROTO_IN}}

# build API (native)
build: proto
	go build -o {{API_EXECUTABLE}}

# run API locally
run_dev: build
	eval $(cat .env.development) {{API_EXECUTABLE}}

# build API (docker) 
build_image:
	{{DOCKER_CMD}} build --build-arg "VERSION={{API_DOCKER_IMAGE}}" -t bjack-api .

# run API locally (docker)
run_image_dev: build_image
	{{DOCKER_CMD}} run -p 8000:8000 --env-file .env.development bjack-api

test: proto
	go test ./...

lint: proto
	golangci-lint run

fmt: proto
	test -z $(gofmt -l .) || gofmt -l . | false

fmt_fix: proto
	gofmt -s -w .

