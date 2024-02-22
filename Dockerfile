ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk --no-cache --virtual .build-deps add make gcc musl-dev binutils-gold git

COPY . /app
WORKDIR /app

ARG TOKEN
ENV TOKEN=${TOKEN}

RUN make build

FROM alpine:${ALPINE_VERSION}

ARG git_commit=unspecified
LABEL git_commit="${git_commit}"
LABEL maintainer="jvolak@frinx.io"
LABEL org.opencontainers.image.source="https://github.com/FRINXio/krakend-ce"

RUN apk add --no-cache ca-certificates tzdata curl && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3, "plugin": { "pattern": ".so", "folder": "/usr/local/lib/krakend/" }}' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend
COPY --from=builder /app/*.so /usr/local/lib/krakend/

RUN chown -R krakend /etc/ssl/certs
USER krakend

COPY startup.sh /startup.sh

WORKDIR /etc/krakend

ENTRYPOINT [ "/startup.sh" ]

EXPOSE 8000 8090
