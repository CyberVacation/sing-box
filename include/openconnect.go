//go:build with_openconnect

package include

import (
	"github.com/CyberVacation/rostra/adapter/endpoint"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/protocol/openconnect"
)

func registerOpenConnectEndpoint(registry *endpoint.Registry) {
	openconnect.RegisterEndpoint(registry)
}

func registerOpenConnectDNSTransport(registry *dns.TransportRegistry) {
	openconnect.RegisterDNSTransport(registry)
}
