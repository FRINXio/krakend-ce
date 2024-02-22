# Krakend OAuth2 Proxy Plugin

This plugin integrates OAuth2 proxy functionality with KrakenD, allowing handling RBAC in Frinx Machine backend services.

## Description

The plugin maps headers from the OAuth2 proxy to headers compatible with Frinx Machine. This mapping is done on request headers. The headers involved in mapping are:

- `X-Tenant-ID`: Header for tenant identification. Configured by by `OAUTH2_KRAKEND_PLUGIN_TENANT_ID`
- `X-Auth-User-Roles`: Frinx header for user roles. Value copied from request header configured by `OAUTH2_KRAKEND_PLUGIN_USER_ROLES_MAP`
- `X-Auth-User-Groups`: Frinx Header for user groups. Value copied from request header configured by `OAUTH2_KRAKEND_PLUGIN_USER_GROUPS_MAP`
- `From`: Header specifying the source of the request. Value copied from request header configured by `OAUTH2_KRAKEND_PLUGIN_FROM_MAP`

Ensure that your OAuth2 proxy and KrakenD configurations are aligned with the usage of these headers for seamless integration.

## Configuration

To use the Krakend OAuth2 Proxy Plugin, follow these steps:

1. **Integrate with your KrakenD Configuration**: Include the plugin in your KrakenD configuration file.

2. **Set Environment Variables**: Ensure that the necessary environment variables are properly set for the plugin to function. These variables include:
   - `OAUTH2_KRAKEND_PLUGIN_TENANT_ID`: Specifies the OAuth2 tenant ID. Default value is `"frinx"`.
   - `OAUTH2_KRAKEND_PLUGIN_USER_ROLES_MAP`: Specifies User-Role request header name. Default value is `"X-Forwarded-Roles"`.
   - `OAUTH2_KRAKEND_PLUGIN_USER_GROUPS_MAP`: Specifies User-Group request header name. Default value is `"X-Forwarded-Groups"`.
   - `OAUTH2_KRAKEND_PLUGIN_FROM_MAP`: Specifies User Identity request header name. Default value is `"X-Forwarded-User"`.

3. **Start your KrakenD Server**: Start your KrakenD server with the OAuth2 proxy plugin integrated.

```json
{
    "version": 3,
    "plugin": {
        "pattern": ".so",
        "folder": "/usr/local/lib/krakend/"
    },
    "extra_config": {
      "plugin/http-server": {
        "name": ["krakend-oauth2-proxy"]
      }
    },
    "endpoints": [
    ]
}
```