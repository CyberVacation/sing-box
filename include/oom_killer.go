package include

import (
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/service/oomkiller"
)

func registerOOMKillerService(registry *service.Registry) {
	oomkiller.RegisterService(registry)
}
