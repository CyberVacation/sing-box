package v2rayxhttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/tls"
	"github.com/sagernet/sing-box/option"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var (
	_ adapter.V2RayMultiplexClientTransport = (*Client)(nil)
	_ adapter.IdleConnectionKeeper          = (*Client)(nil)
)

type Client struct {
	ctx     context.Context
	config  clientConfig
	pool    *clientPool
	mu      sync.Mutex
	streams map[*streamConn]struct{}
}

func NewClient(ctx context.Context, rawDialer N.Dialer, serverAddr M.Socksaddr, options option.V2RayXHTTPOptions, tlsConfig tls.Config) (*Client, error) {
	config, err := newConfig(options, serverAddr, tlsConfig)
	if err != nil {
		return nil, err
	}
	factory, err := transportFactory(rawDialer, serverAddr, config.version, tlsConfig, config.pool.keepAlive)
	if err != nil {
		return nil, err
	}
	return &Client{
		ctx:     ctx,
		config:  config,
		pool:    newClientPool(config.pool, factory),
		streams: make(map[*streamConn]struct{}),
	}, nil
}

func (c *Client) DialContext(ctx context.Context) (net.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := c.ctx.Err(); err != nil {
		return nil, err
	}
	session := ""
	if c.config.mode != modeSingle {
		id, err := c.config.metadata.newSessionID()
		if err != nil {
			return nil, err
		}
		session = id
	}
	conn, streamCtx := newStream(ctx, c.config.server)
	// A stream owns its cancellation; transport Close resets all existing streams.
	c.mu.Lock()
	uploadLease := c.pool.acquire()
	c.streams[conn] = struct{}{}
	conn.onClose = func() {
		c.mu.Lock()
		delete(c.streams, conn)
		c.mu.Unlock()
		uploadLease.Close()
	}
	c.mu.Unlock()
	stop := context.AfterFunc(c.ctx, func() {
		conn.finish(c.ctx.Err())
	})
	go func() {
		<-streamCtx.Done()
		stop()
		conn.finish(streamCtx.Err())
	}()
	if c.config.mode == modeSingle {
		go c.stream(streamCtx, conn, uploadLease, session, conn.wire, false)
	} else {
		go c.stream(streamCtx, conn, uploadLease, session, nil, false)
		if c.config.mode == modeStream {
			go c.stream(streamCtx, conn, uploadLease, session, conn.wire, true)
		} else {
			go c.uploadPackets(streamCtx, conn, uploadLease, session)
		}
	}
	return conn, nil
}

func (c *Client) stream(ctx context.Context, conn *streamConn, lease *poolLease, session string, body io.Reader, upload bool) {
	req := c.config.newRequest(ctx, session, "", body, nil)
	response, err := lease.RoundTrip(req)
	if err != nil {
		conn.finish(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		conn.finish(fmt.Errorf("xhttp: unexpected HTTP status: %s", response.Status))
		return
	}
	var destination io.Writer = conn.wire
	if upload {
		destination = io.Discard
	}
	_, err = io.Copy(destination, response.Body)
	if upload && err == nil {
		// The upload response can end before the independent download has
		// drained. Let the download close the connection after its final bytes.
		return
	}
	if err == nil {
		err = io.EOF
	}
	conn.finish(err)
}

func (c *Client) MultiplexEnabled() bool {
	return c.config.version != "1.1"
}

func (c *Client) SetKeepIdleConnections(keep bool) {
	c.pool.setKeepIdle(keep)
}

func (c *Client) CloseIdleConnections() {
	c.pool.closeIdle()
}

func (c *Client) Close() error {
	c.mu.Lock()
	streams := make([]*streamConn, 0, len(c.streams))
	for s := range c.streams {
		streams = append(streams, s)
	}
	c.mu.Unlock()
	for _, s := range streams {
		s.Close()
	}
	c.pool.reset()
	return nil
}
