//go:build !with_ocm

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

func registerOCMService(registry *service.Registry) {
	service.Register[option.OCMServiceOptions](registry, C.TypeOCM, func(ctx context.Context, logger log.ContextLogger, tag string, options option.OCMServiceOptions) (adapter.Service, error) {
		return nil, E.New(`OCM is not included in this build, rebuild with -tags with_ocm`)
	})
}
