#!/bin/sh
cd krakend-azure-plugin/
go build -trimpath -buildmode=plugin -o azure_plugin.so
cp azure_plugin.so ../
cd ..
