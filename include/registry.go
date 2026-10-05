package include

import (
	"context"

	"github.com/CyberVacation/rostra"
	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/adapter/certificate"
	"github.com/CyberVacation/rostra/adapter/endpoint"
	"github.com/CyberVacation/rostra/adapter/inbound"
	"github.com/CyberVacation/rostra/adapter/outbound"
	"github.com/CyberVacation/rostra/adapter/service"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/dns"
	"github.com/CyberVacation/rostra/dns/transport"
	"github.com/CyberVacation/rostra/dns/transport/fakeip"
	"github.com/CyberVacation/rostra/dns/transport/hosts"
	"github.com/CyberVacation/rostra/dns/transport/local"
	"github.com/CyberVacation/rostra/dns/transport/mdns"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"
	"github.com/CyberVacation/rostra/protocol/anytls"
	"github.com/CyberVacation/rostra/protocol/block"
	"github.com/CyberVacation/rostra/protocol/bridge"
	"github.com/CyberVacation/rostra/protocol/direct"
	"github.com/CyberVacation/rostra/protocol/group"
	"github.com/CyberVacation/rostra/protocol/http"
	"github.com/CyberVacation/rostra/protocol/masque"
	"github.com/CyberVacation/rostra/protocol/mixed"
	"github.com/CyberVacation/rostra/protocol/naive"
	"github.com/CyberVacation/rostra/protocol/redirect"
	"github.com/CyberVacation/rostra/protocol/shadowsocks"
	"github.com/CyberVacation/rostra/protocol/shadowtls"
	"github.com/CyberVacation/rostra/protocol/snell"
	"github.com/CyberVacation/rostra/protocol/socks"
	"github.com/CyberVacation/rostra/protocol/ssh"
	"github.com/CyberVacation/rostra/protocol/tor"
	"github.com/CyberVacation/rostra/protocol/trojan"
	"github.com/CyberVacation/rostra/protocol/tun"
	"github.com/CyberVacation/rostra/protocol/vless"
	"github.com/CyberVacation/rostra/protocol/vmess"
	"github.com/CyberVacation/rostra/service/api"
	originca "github.com/CyberVacation/rostra/service/origin_ca"
	"github.com/CyberVacation/rostra/service/resolved"
	"github.com/CyberVacation/rostra/service/ssmapi"

	E "github.com/sagernet/sing/common/exceptions"
)

func Context(ctx context.Context) context.Context {
	return box.Context(ctx, InboundRegistry(), OutboundRegistry(), EndpointRegistry(), DNSTransportRegistry(), ServiceRegistry(), CertificateProviderRegistry())
}

func InboundRegistry() *inbound.Registry {
	registry := inbound.NewRegistry()

	tun.RegisterInbound(registry)
	redirect.RegisterRedirect(registry)
	redirect.RegisterTProxy(registry)
	direct.RegisterInbound(registry)

	socks.RegisterInbound(registry)
	http.RegisterInbound(registry)
	mixed.RegisterInbound(registry)

	shadowsocks.RegisterInbound(registry)
	snell.RegisterInbound(registry)
	vmess.RegisterInbound(registry)
	trojan.RegisterInbound(registry)
	naive.RegisterInbound(registry)
	shadowtls.RegisterInbound(registry)
	vless.RegisterInbound(registry)
	anytls.RegisterInbound(registry)

	registerQUICInbounds(registry)
	registerCloudflaredInbound(registry)
	registerTailcatInbound(registry)
	registerStubForRemovedInbounds(registry)

	return registry
}

func OutboundRegistry() *outbound.Registry {
	registry := outbound.NewRegistry()

	direct.RegisterOutbound(registry)
	bridge.RegisterOutbound(registry)

	block.RegisterOutbound(registry)

	group.RegisterSelector(registry)
	group.RegisterURLTest(registry)
	group.RegisterLoadBalance(registry)

	socks.RegisterOutbound(registry)
	http.RegisterOutbound(registry)
	shadowsocks.RegisterOutbound(registry)
	snell.RegisterOutbound(registry)
	vmess.RegisterOutbound(registry)
	trojan.RegisterOutbound(registry)
	registerNaiveOutbound(registry)
	tor.RegisterOutbound(registry)
	ssh.RegisterOutbound(registry)
	shadowtls.RegisterOutbound(registry)
	vless.RegisterOutbound(registry)
	anytls.RegisterOutbound(registry)

	registerQUICOutbounds(registry)
	registerTailcatOutbound(registry)
	registerStubForRemovedOutbounds(registry)

	return registry
}

func EndpointRegistry() *endpoint.Registry {
	registry := endpoint.NewRegistry()

	registerWireGuardEndpoint(registry)
	registerOpenConnectEndpoint(registry)
	registerOpenVPNEndpoints(registry)
	masque.RegisterEndpoint(registry)
	registerTailscaleEndpoint(registry)

	return registry
}

func DNSTransportRegistry() *dns.TransportRegistry {
	registry := dns.NewTransportRegistry()

	transport.RegisterTCP(registry)
	transport.RegisterUDP(registry)
	transport.RegisterTLS(registry)
	transport.RegisterHTTPS(registry)
	hosts.RegisterTransport(registry)
	local.RegisterTransport(registry)
	mdns.RegisterTransport(registry)
	fakeip.RegisterTransport(registry)
	resolved.RegisterTransport(registry)

	registerQUICTransports(registry)
	registerDHCPTransport(registry)
	registerTailscaleTransport(registry)
	registerOpenConnectDNSTransport(registry)
	registerOpenVPNDNSTransport(registry)

	return registry
}

func ServiceRegistry() *service.Registry {
	registry := service.NewRegistry()

	api.RegisterService(registry)
	resolved.RegisterService(registry)
	ssmapi.RegisterService(registry)

	registerQUICServices(registry)
	registerDERPService(registry)
	registerCCMService(registry)
	registerOCMService(registry)
	registerOOMKillerService(registry)
	registerUSBIPServices(registry)

	return registry
}

func CertificateProviderRegistry() *certificate.Registry {
	registry := certificate.NewRegistry()

	registerACMECertificateProvider(registry)
	registerTailscaleCertificateProvider(registry)
	originca.RegisterCertificateProvider(registry)

	return registry
}

func registerStubForRemovedInbounds(registry *inbound.Registry) {
	inbound.Register[option.ShadowsocksInboundOptions](registry, C.TypeShadowsocksR, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowsocksInboundOptions) (adapter.Inbound, error) {
		return nil, E.New("ShadowsocksR is deprecated and removed in sing-box 1.6.0")
	})
}

func registerStubForRemovedOutbounds(registry *outbound.Registry) {
	outbound.Register[option.ShadowsocksROutboundOptions](registry, C.TypeShadowsocksR, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowsocksROutboundOptions) (adapter.Outbound, error) {
		return nil, E.New("ShadowsocksR is deprecated and removed in sing-box 1.6.0")
	})
	outbound.Register[option.StubOptions](registry, C.TypeWireGuard, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.StubOptions) (adapter.Outbound, error) {
		return nil, E.New("WireGuard outbound is deprecated in sing-box 1.11.0 and removed in sing-box 1.13.0, use WireGuard endpoint instead")
	})
}
