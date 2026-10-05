//go:build !with_ccm

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/adapter/service"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func registerCCMService(registry *service.Registry) {
	service.Register[option.CCMServiceOptions](registry, C.TypeCCM, func(ctx context.Context, logger log.ContextLogger, tag string, options option.CCMServiceOptions) (adapter.Service, error) {
		return nil, E.New(`CCM is not included in this build, rebuild with -tags with_CCM`)
	})
}
