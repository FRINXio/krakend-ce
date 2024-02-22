// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

var pluginName = "krakend-oauth2-proxy"
var HandlerRegisterer = registerer(pluginName)

type registerer string

var oauthTenantHeader string
var oauthRolesHeaderMap string
var oauthGroupsHeaderMap string
var oauthFromHeaderMap string

var frinxTenantHeader string = "X-Tenant-ID"
var frinxRolesHeader string = "X-Auth-User-Roles"
var frinxGroupsHeader string = "X-Auth-User-Groups"
var frinxFromHeaderMap string = "From"

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Remove incoming headers if exists
		// Prevent injecting RBAC headers in request before request comes to KrakenD
		req.Header.Del(frinxTenantHeader)
		req.Header.Del(frinxRolesHeader)
		req.Header.Del(frinxGroupsHeader)
		req.Header.Del(frinxFromHeaderMap)

		// Set static tenant id heder from ENV variables
		// Tenant is not forwarded from oauth2-proxy,
		// but Frinx Machine backend services depends on this header
		req.Header.Add(frinxTenantHeader, oauthTenantHeader)
		logger.Debug(fmt.Sprintf("Authorization header map: %s - %s", frinxTenantHeader, req.Header.Get(frinxTenantHeader)))

		// Map headers from oauth2-proxy to Frinx Machine compatible headers
		if req.Header.Get(oauthRolesHeaderMap) != "" {
			req.Header.Add(frinxRolesHeader, req.Header.Get(oauthRolesHeaderMap))
			logger.Debug(fmt.Sprintf("Authorization header map: %s - %s", frinxRolesHeader, req.Header.Get(frinxRolesHeader)))
		}

		// Map headers from oauth2-proxy to Frinx Machine compatible headers
		if req.Header.Get(oauthGroupsHeaderMap) != "" {
			req.Header.Add(frinxGroupsHeader, req.Header.Get(oauthGroupsHeaderMap))
			logger.Debug(fmt.Sprintf("Authorization header map: %s - %s", frinxGroupsHeader, req.Header.Get(frinxGroupsHeader)))
		}

		// Map headers from oauth2-proxy to Frinx Machine compatible headers
		if req.Header.Get(oauthFromHeaderMap) != "" {
			req.Header.Add(frinxFromHeaderMap, req.Header.Get(oauthFromHeaderMap))
			logger.Debug(fmt.Sprintf("Authorization header map: %s - %s", frinxFromHeaderMap, req.Header.Get(frinxFromHeaderMap)))
		}

		handler.ServeHTTP(w, req)

	}), nil
}

func init() {

	oauthTenantHeader = os.Getenv("OAUTH2_KRAKEND_PLUGIN_TENANT_ID")
	oauthRolesHeaderMap = os.Getenv("OAUTH2_KRAKEND_PLUGIN_USER_ROLES_MAP")
	oauthGroupsHeaderMap = os.Getenv("OAUTH2_KRAKEND_PLUGIN_USER_GROUPS_MAP")
	oauthFromHeaderMap = os.Getenv("OAUTH2_KRAKEND_PLUGIN_FROM_MAP")

	if oauthTenantHeader == "" {
		oauthTenantHeader = "frinx"
		logger.Warning(fmt.Sprintf("WARN: no OAUTH2_KRAKEND_PLUGIN_TENANT_ID, using default: %s \n", oauthTenantHeader))
	}

	if oauthRolesHeaderMap == "" {
		oauthRolesHeaderMap = "X-Forwarded-Roles"
		logger.Warning(fmt.Sprintf("WARN: OAUTH2_KRAKEND_PLUGIN_USER_ROLES_MAP set, using default: %s \n", oauthRolesHeaderMap))
	}

	if oauthGroupsHeaderMap == "" {
		oauthGroupsHeaderMap = "X-Forwarded-Groups"
		logger.Warning(fmt.Sprintf("WARN: OAUTH2_KRAKEND_PLUGIN_USER_GROUPS_MAP set, using default: %s \n", oauthGroupsHeaderMap))
	}

	if oauthFromHeaderMap == "" {
		oauthFromHeaderMap = "X-Forwarded-User"
		logger.Warning(fmt.Sprintf("WARN: OAUTH2_KRAKEND_PLUGIN_FROM_MAP set, using default: %s \n", oauthFromHeaderMap))
	}
}

func main() {}

var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Info(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
