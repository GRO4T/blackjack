# TODO: Rename blackjack to bjackapi
SERVER_DIR=./backend
PROTO_IN=./blackjack.proto
PROTO_OUT_DIR=./backend/proto
PROTO_OUT=./backend/proto/blackjack.pb.go ./backend/proto/blackjack_grpc.pb.go
EXECUTABLE=blackjack

all: proto build

build:
	cd $(SERVER_DIR) && go build -o ./bin/$(EXECUTABLE)

install:
	cd $(SERVER_DIR) && go install

test:
	cd $(SERVER_DIR) && go test ./...

lint:
	cd $(SERVER_DIR) && golangci-lint run

clean:
	rm -f $(PROTO_OUT) $(SERVER_DIR)/bin/$(EXECUTABLE)

proto: $(PROTO_OUT)

$(PROTO_OUT): $(PROTO_IN)
	protoc -I=. --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_IN)