//go:build with_wireguard

package include

import (
	"github.com/CyberVacation/rostra/adapter/endpoint"
	"github.com/CyberVacation/rostra/protocol/wireguard"
)

func registerWireGuardEndpoint(registry *endpoint.Registry) {
	wireguard.RegisterEndpoint(registry)
}
