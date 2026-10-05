package sniff_test

import (
	"context"
	"encoding/hex"
	"os"
	"testing"

	"github.com/CyberVacation/rostra/adapter"
	"github.com/CyberVacation/rostra/common/sniff"
	C "github.com/CyberVacation/rostra/constant"

	"github.com/stretchr/testify/require"
)

func TestSniffNTP(t *testing.T) {
	t.Parallel()
	packet, err := hex.DecodeString("1b0006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	var metadata adapter.InboundContext
	err = sniff.NTP(context.Background(), &metadata, packet)
	require.NoError(t, err)
	require.Equal(t, metadata.Protocol, C.ProtocolNTP)
}

func TestSniffNTPFailed(t *testing.T) {
	t.Parallel()
	packet, err := hex.DecodeString("400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	var metadata adapter.InboundContext
	err = sniff.NTP(context.Background(), &metadata, packet)
	require.ErrorIs(t, err, os.ErrInvalid)
}
