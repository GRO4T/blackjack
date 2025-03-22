FROM golang:1.23.4-alpine3.21

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.62.2 \
    && apk update \
    && apk add --no-cache git just \
    && rm -rf /var/cache/apk/* \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v29.2/protoc-29.2-linux-x86_64.zip \
    && unzip protoc-29.2-linux-x86_64.zip -d $HOME/.local \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
