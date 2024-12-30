# TODO: Rename blackjack to bjackapi
SERVER_DIR=./backend
PROTO_IN=./blackjack.proto
PROTO_OUT_DIR=./backend/grpc
PROTO_OUT=./backend/grpc/blackjack.pb.go ./backend/grpc/blackjack_grpc.pb.go
EXECUTABLE=./backend/blackjack

all: proto build

build:
	cd $(SERVER_DIR) && go build

install:
	cd $(SERVER_DIR) && go install

clean:
	rm -f $(PROTO_OUT) $(EXECUTABLE)

proto: $(PROTO_OUT)

$(PROTO_OUT): $(PROTO_IN)
	protoc -I=. --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_IN)