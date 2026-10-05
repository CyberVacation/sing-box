//go:build !with_openconnect

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/adapter/endpoint"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func registerOpenConnectEndpoint(registry *endpoint.Registry) {
	endpoint.Register[option.OpenConnectEndpointOptions](registry, C.TypeOpenConnect, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.OpenConnectEndpointOptions) (adapter.Endpoint, error) {
		return nil, E.New(`OpenConnect is not included in this build, rebuild with -tags with_openconnect`)
	})
}

func registerOpenConnectDNSTransport(registry *dns.TransportRegistry) {
	dns.RegisterTransport[option.OpenConnectDNSServerOptions](registry, C.DNSTypeOpenConnect, func(ctx context.Context, logger log.ContextLogger, tag string, options option.OpenConnectDNSServerOptions) (adapter.DNSTransport, error) {
		return nil, E.New(`OpenConnect is not included in this build, rebuild with -tags with_openconnect`)
	})
}
