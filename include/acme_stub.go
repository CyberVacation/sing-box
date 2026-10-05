//go:build !with_acme

package include

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/adapter/certificate"
	C "github.com/CyberVacation/rostra/constant"
	"github.com/CyberVacation/rostra/log"
	"github.com/CyberVacation/rostra/option"

	E "github.com/sagernet/sing/common/exceptions"
)

func registerACMECertificateProvider(registry *certificate.Registry) {
	certificate.Register[option.ACMECertificateProviderOptions](registry, C.TypeACME, func(ctx context.Context, logger log.ContextLogger, tag string, options option.ACMECertificateProviderOptions) (adapter.CertificateProviderService, error) {
		return nil, E.New(`ACME is not included in this build, rebuild with -tags with_acme`)
	})
}
