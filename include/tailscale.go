//go:build with_tailscale

package include

import (
	"github.com/CyberVacation/rostra/adapter/certificate"
	"github.com/CyberVacation/rostra/adapter/endpoint"
	"github.com/CyberVacation/rostra/adapter/inbound"
	"github.com/CyberVacation/rostra/adapter/outbound"
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/protocol/tailscale"
	"github.com/CyberVacation/rostra/service/derp"
)

func registerTailscaleEndpoint(registry *endpoint.Registry) {
	tailscale.RegisterEndpoint(registry)
}

func registerTailcatInbound(registry *inbound.Registry) {
	tailscale.RegisterTailcatInbound(registry)
}

func registerTailcatOutbound(registry *outbound.Registry) {
	tailscale.RegisterTailcatOutbound(registry)
}

func registerTailscaleTransport(registry *dns.TransportRegistry) {
	tailscale.RegistryTransport(registry)
}

func registerTailscaleCertificateProvider(registry *certificate.Registry) {
	tailscale.RegisterCertificateProvider(registry)
}

func registerDERPService(registry *service.Registry) {
	derp.Register(registry)
}
