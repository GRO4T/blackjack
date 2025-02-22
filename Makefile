SERVER_DIR=./bjack-api
PROTO_IN=./blackjack.proto
PROTO_OUT_DIR=./bjack-api/proto
PROTO_OUT=./bjack-api/proto/blackjack.pb.go ./bjack-api/proto/blackjack_grpc.pb.go
EXECUTABLE=bjack-api

all: install

setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	cd bjack-ui && npm install

build: proto
	cd $(SERVER_DIR) && go build -o ./bin/$(EXECUTABLE)

install: proto
	cd $(SERVER_DIR) && go install

test: proto
	cd $(SERVER_DIR) && go test ./...

lint: proto
	cd $(SERVER_DIR) && golangci-lint run

fmt: proto
	cd $(SERVER_DIR) && gofmt -s -w .

clean:
	rm -f $(PROTO_OUT) $(SERVER_DIR)/bin/$(EXECUTABLE)

proto: $(PROTO_OUT)

$(PROTO_OUT): $(PROTO_IN)
	protoc -I=. --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_IN)

ci_image:
	docker build -t dkolaska/blackjack-ci:0.0.3 -f ci.dockerfile --platform linux/amd64 .
