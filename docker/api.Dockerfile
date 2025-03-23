FROM golang:1.23.4-alpine3.21 AS build

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN apk update \
    && apk add --no-cache just \
    && rm -rf /var/cache/apk/* \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v29.2/protoc-29.2-linux-x86_64.zip \
    && unzip protoc-29.2-linux-x86_64.zip -d $HOME/.local

COPY justfile /
RUN just setup_api
COPY blackjack.proto /
COPY bjack-api/ /bjack-api
RUN just build_api 

FROM golang:1.23.4-alpine3.21
COPY --from=build /bjack-api/bin/bjack-api /bjack-api
EXPOSE 8000
ENTRYPOINT "/bjack-api"
