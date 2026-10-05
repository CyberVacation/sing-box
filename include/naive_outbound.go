//go:build with_naive_outbound

package include

import (
	"github.com/CyberVacation/rostra/adapter/outbound"
	"github.com/CyberVacation/rostra/protocol/naive"
)

func registerNaiveOutbound(registry *outbound.Registry) {
	naive.RegisterOutbound(registry)
}
