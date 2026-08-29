ARG VERSION

FROM ${VERSION} AS build

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN apk update \
    && apk add --no-cache just \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY justfile api.just /app
RUN just api setup


COPY go.mod go.sum /app
RUN go mod download

COPY . /app

RUN just api test
RUN just api fmt
RUN just api lint
RUN just api build

FROM ${VERSION}
COPY --from=build /app/bin/bjack-api /bjack-api
EXPOSE 8000
ENTRYPOINT "/bjack-api"

