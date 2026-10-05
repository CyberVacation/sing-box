//go:build !with_quic

package networkquality

import (
	C "github.com/CyberVacation/rostra/constant"

	N "github.com/sagernet/sing/common/network"
)

func NewHTTP3MeasurementClientFactory(dialer N.Dialer) (MeasurementClientFactory, error) {
	return nil, C.ErrQUICNotIncluded
}
