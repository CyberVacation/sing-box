//go:build with_dhcp

package include

import (
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/dns/transport/dhcp"
)

func registerDHCPTransport(registry *dns.TransportRegistry) {
	dhcp.RegisterTransport(registry)
}
