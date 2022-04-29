FROM debian:stable-slim

LABEL maintainer="jvolak@frinx.io"

RUN apt-get update && apt-get upgrade -y && \
	apt-get install -y ca-certificates curl && \
	rm -rf /var/lib/apt/lists/*

ADD krakend /usr/bin/krakend
ADD azure_plugin.so /usr/local/lib/krakend/azure_plugin.so
ADD startup.sh /startup.sh

RUN useradd -r -c "KrakenD user" -U krakend

RUN chown -R krakend /etc/ssl/certs

USER krakend

VOLUME [ "/etc/krakend" ]

WORKDIR /etc/krakend

ENTRYPOINT [ "/startup.sh" ]

EXPOSE 8000 8090
