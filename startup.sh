#!/bin/sh
set -a

update-ca-certificates

if [ -f "/set_env_secrets.sh" ]; then
  . /set_env_secrets.sh ''
fi

if [ "${PROXY_ENABLED}" != 'true' ]; then
  unset HTTP_PROXY HTTPS_PROXY NO_PROXY
fi

/usr/bin/krakend run -c /etc/krakend/krakend.json
