//go:build with_quic

package include

import (
	"github.com/CyberVacation/rostra/adapter/inbound"
	"github.com/CyberVacation/rostra/adapter/outbound"
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/dns/transport/quic"
	"github.com/CyberVacation/rostra/protocol/hysteria"
	"github.com/CyberVacation/rostra/protocol/hysteria2"
	_ "github.com/CyberVacation/rostra/protocol/naive/quic"
	"github.com/CyberVacation/rostra/protocol/tuic"
	_ "github.com/CyberVacation/rostra/transport/v2rayquic"
)

func registerQUICInbounds(registry *inbound.Registry) {
	hysteria.RegisterInbound(registry)
	tuic.RegisterInbound(registry)
	hysteria2.RegisterInbound(registry)
}

func registerQUICOutbounds(registry *outbound.Registry) {
	hysteria.RegisterOutbound(registry)
	tuic.RegisterOutbound(registry)
	hysteria2.RegisterOutbound(registry)
}

func registerQUICTransports(registry *dns.TransportRegistry) {
	quic.RegisterTransport(registry)
	quic.RegisterHTTP3Transport(registry)
}

func registerQUICServices(registry *service.Registry) {
	hysteria2.RegisterRealmService(registry)
}
