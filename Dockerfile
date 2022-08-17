ARG GOLANG_VERSION
ARG ALPINE_VERSION

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk --no-cache add make gcc musl-dev binutils-gold git

COPY . /app
WORKDIR /app

ARG GH_TOKEN
ENV GH_TOKEN=${GH_TOKEN}

RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/FRINXio".insteadOf "https://github.com/FRINXio" && make build

FROM alpine:${ALPINE_VERSION}

LABEL maintainer="jvolak@frinx.io"

RUN apk add --no-cache ca-certificates curl && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend

COPY azure_plugin.so /usr/local/lib/krakend/azure_plugin.so

RUN chown -R krakend /etc/ssl/certs
USER krakend

COPY startup.sh /startup.sh

WORKDIR /etc/krakend

ENTRYPOINT [ "/startup.sh" ]

EXPOSE 8000 8090
