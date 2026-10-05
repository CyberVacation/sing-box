//go:build !with_naive_outbound

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/adapter/outbound"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func registerNaiveOutbound(registry *outbound.Registry) {
	outbound.Register[option.NaiveOutboundOptions](registry, C.TypeNaive, func(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.NaiveOutboundOptions) (adapter.Outbound, error) {
		return nil, E.New(`naive outbound is not included in this build, rebuild with -tags with_naive_outbound`)
	})
}
