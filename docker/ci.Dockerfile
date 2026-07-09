ARG VERSION
FROM ${VERSION}

ENV HOME /root
ENV PATH="$PATH:$HOME/.local/bin"

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.62.2 \
    && apk update \
    && apk add --no-cache git just npm \
    && rm -rf /var/cache/apk/* 

