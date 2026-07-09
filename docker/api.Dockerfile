ARG VERSION

FROM ${VERSION} AS build

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN apk update \
    && apk add --no-cache just \
    && rm -rf /var/cache/apk/*

COPY justfile /
RUN just setup_api
COPY blackjack.proto /
COPY bjack-api/ /bjack-api
RUN just build_api 

FROM ${VERSION}
COPY --from=build /bjack-api/bin/bjack-api /bjack-api
EXPOSE 8000
ENTRYPOINT "/bjack-api"
