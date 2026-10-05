//go:build with_usbip && (linux || (darwin && cgo) || windows)

package include

import (
	"github.com/CyberVacation/rostra/adapter/service"
	"github.com/CyberVacation/rostra/service/usbip"
)

func registerUSBIPServices(registry *service.Registry) {
	usbip.RegisterService(registry)
}
