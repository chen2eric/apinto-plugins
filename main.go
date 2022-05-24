package main

import (
	circuit_breaker "github.com/eolinker/apinto-plugins/circuit-breaker"
	"github.com/eolinker/apinto-plugins/cors"
	extra_params "github.com/eolinker/apinto-plugins/extra-params"
	"github.com/eolinker/apinto-plugins/gzip"
	ip_restriction "github.com/eolinker/apinto-plugins/ip-restriction"
	params_transformer "github.com/eolinker/apinto-plugins/params-transformer"
	proxy_rewrite "github.com/eolinker/apinto-plugins/proxy-rewrite"
	rate_limiting "github.com/eolinker/apinto-plugins/rate-limiting"
	response_rewrite "github.com/eolinker/apinto-plugins/response-rewrite"
	"github.com/eolinker/eosc"
)

type builder struct {
}

func (b *builder) Register(register eosc.IExtenderDriverRegister) {
	ip_restriction.Register(register)
	cors.Register(register)
	circuit_breaker.Register(register)
	extra_params.Register(register)
	gzip.Register(register)
	params_transformer.Register(register)
	proxy_rewrite.Register(register)
	rate_limiting.Register(register)
	response_rewrite.Register(register)
}

func Builder() eosc.ExtenderBuilder {
	return new(builder)
}
