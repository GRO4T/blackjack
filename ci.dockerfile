FROM golang:1.23.4-alpine3.21

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.62.2 \
    && apk update \
    && apk add --no-cache git just npm \
    && rm -rf /var/cache/apk/* \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v29.2/protoc-29.2-linux-x86_64.zip \
    && unzip protoc-29.2-linux-x86_64.zip -d $HOME/.local
