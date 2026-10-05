//go:build with_acme

package include

import (
	"github.com/CyberVacation/rostra/adapter/certificate"
	"github.com/CyberVacation/rostra/service/acme"
)

func registerACMECertificateProvider(registry *certificate.Registry) {
	acme.RegisterCertificateProvider(registry)
}
