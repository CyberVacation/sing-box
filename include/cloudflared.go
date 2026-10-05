//go:build with_cloudflared

package include

import (
	"github.com/CyberVacation/rostra/adapter/inbound"
	"github.com/CyberVacation/rostra/protocol/cloudflare"
)

func registerCloudflaredInbound(registry *inbound.Registry) {
	cloudflare.RegisterInbound(registry)
}
