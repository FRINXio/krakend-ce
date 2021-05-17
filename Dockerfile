FROM debian:buster-slim

LABEL maintainer="jvolak@frinx.io"

ADD certs /usr/local/share/ca-certificates

RUN apt-get update && \
	apt-get install -y ca-certificates curl && \
	update-ca-certificates && \
	rm -rf /var/lib/apt/lists/*

ADD krakend /usr/bin/krakend
ADD azure_plugin.so /usr/local/lib/krakend/azure_plugin.so

RUN useradd -r -c "KrakenD user" -U krakend

USER krakend

VOLUME [ "/etc/krakend" ]

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
