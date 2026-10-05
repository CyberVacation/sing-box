//go:build !with_clash_api

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/experimental"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func init() {
	experimental.RegisterClashServerConstructor(func(ctx context.Context, logFactory log.ObservableFactory, options option.ClashAPIOptions) (adapter.LifecycleService, error) {
		return nil, E.New(`clash api is not included in this build, rebuild with -tags with_clash_api`)
	})
}
