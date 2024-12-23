FROM golang:1.23.4-alpine3.21

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.62.2