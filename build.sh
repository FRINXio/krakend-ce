#!/bin/sh
cd krakend-azure-plugin/
go build -trimpath -buildmode=plugin -buildvcs=false -o azure_plugin.so
cp azure_plugin.so ../
cd ..

cd krakend-oauth2-proxy-plugin/
go build -trimpath -buildmode=plugin -o oauth2_proxy_plugin.so .
cp oauth2_proxy_plugin.so ../
cd ..
