//go:build with_openvpn

package include

import (
	"github.com/CyberVacation/rostra/adapter/endpoint"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/protocol/openvpn"
)

func registerOpenVPNEndpoints(registry *endpoint.Registry) {
	openvpn.RegisterEndpoint(registry)
}

func registerOpenVPNDNSTransport(registry *dns.TransportRegistry) {
	openvpn.RegisterDNSTransport(registry)
}
