#!/bin/sh
set -a

update-ca-certificates

if [ -f "/set_env_secrets.sh" ]; then
  . /set_env_secrets.sh ''
fi

/usr/bin/krakend run -c /etc/krakend/krakend.json
