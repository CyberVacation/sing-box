package v2rayxhttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sagernet/sing-box/option"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/stretchr/testify/require"
)

func TestStreamUploadResponseLifecycle(t *testing.T) {
	uploadError := errors.New("upload response failed")
	for _, test := range []struct {
		name string
		err  error
	}{
		{name: "clean upload completion preserves download"},
		{name: "upload failure closes connection", err: uploadError},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client, err := NewClient(ctx, N.SystemDialer, M.ParseSocksaddr("127.0.0.1:80"), option.V2RayXHTTPOptions{Mode: modeStream, HTTPVersion: "2"}, nil)
			require.NoError(t, err)
			defer client.Close()

			payload := bytes.Repeat([]byte("buffered download data"), 4096)
			uploadDone := make(chan struct{})
			client.pool.factory = func() http.RoundTripper {
				return testRoundTripper(func(request *http.Request) (*http.Response, error) {
					body := io.NopCloser(bytes.NewReader(payload))
					if request.Method == http.MethodPost {
						body = &uploadResponseBody{err: test.err, closed: uploadDone}
					}
					return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
				})
			}
			conn, err := client.DialContext(ctx)
			require.NoError(t, err)
			defer conn.Close()
			// Do not read until the upload response has been fully handled. The
			// independent download can still have unread data at that point.
			select {
			case <-uploadDone:
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			content, err := io.ReadAll(conn)
			if test.err != nil {
				require.ErrorIs(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Len(t, content, len(payload))
				require.Equal(t, payload, content)
			}
		})
	}
}

type testRoundTripper func(*http.Request) (*http.Response, error)

func (f testRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type uploadResponseBody struct {
	err    error
	closed chan struct{}
}

func (b *uploadResponseBody) Read([]byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	return 0, io.EOF
}

func (b *uploadResponseBody) Close() error {
	close(b.closed)
	return nil
}
