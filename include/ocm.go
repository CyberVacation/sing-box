//go:build with_ocm

package include

import (
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/service/ocm"
)

func registerOCMService(registry *service.Registry) {
	ocm.RegisterService(registry)
}
