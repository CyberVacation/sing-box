package settings

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
)

type WIFIMonitor interface {
	ReadWIFIState(ctx context.Context) adapter.WIFIState
	Start() error
	Close() error
}
