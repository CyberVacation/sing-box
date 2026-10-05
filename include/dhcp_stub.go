//go:build !with_dhcp

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func registerDHCPTransport(registry *dns.TransportRegistry) {
	dns.RegisterTransport[option.DHCPDNSServerOptions](registry, C.DNSTypeDHCP, func(ctx context.Context, logger log.ContextLogger, tag string, options option.DHCPDNSServerOptions) (adapter.DNSTransport, error) {
		return nil, E.New(`DHCP is not included in this build, rebuild with -tags with_dhcp`)
	})
}
