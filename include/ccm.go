//go:build with_ccm && (!darwin || cgo)

package include

import (
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/service/ccm"
)

func registerCCMService(registry *service.Registry) {
	ccm.RegisterService(registry)
}
