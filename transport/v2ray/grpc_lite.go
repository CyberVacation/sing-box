//go:build !with_grpc

package v2ray

import (
	"context"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/common/tls"
	"github.com/CyberVacation/rostra/option"
	"github.com/CyberVacation/rostra/transport/v2raygrpclite"

	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

func NewGRPCServer(ctx context.Context, logger logger.ContextLogger, options option.V2RayGRPCOptions, tlsConfig tls.ServerConfig, handler adapter.V2RayServerTransportHandler) (adapter.V2RayServerTransport, error) {
	return v2raygrpclite.NewServer(ctx, logger, options, tlsConfig, handler)
}

func NewGRPCClient(ctx context.Context, dialer N.Dialer, serverAddr M.Socksaddr, options option.V2RayGRPCOptions, tlsConfig tls.Config) (adapter.V2RayClientTransport, error) {
	return v2raygrpclite.NewClient(ctx, dialer, serverAddr, options, tlsConfig), nil
}
