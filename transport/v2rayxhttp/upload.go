package v2rayxhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const maxConcurrentPosts = 8

// uploadBuffer coalesces writes while the sender waits for its next POST slot.
type uploadBuffer struct {
	mu     sync.Mutex
	cond   *sync.Cond
	data   []byte
	limit  int
	closed bool
}

func newUploadBuffer(limit int) *uploadBuffer {
	b := &uploadBuffer{limit: limit}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *uploadBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for len(p) > 0 {
		for len(b.data) == b.limit && !b.closed {
			b.cond.Wait()
		}
		if b.closed {
			return n, io.ErrClosedPipe
		}
		count := min(len(p), b.limit-len(b.data))
		b.data = append(b.data, p[:count]...)
		p = p[count:]
		n += count
		b.cond.Broadcast()
	}
	return n, nil
}

func (b *uploadBuffer) take() ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.data) == 0 && !b.closed {
		b.cond.Wait()
	}
	if b.closed {
		return nil, io.ErrClosedPipe
	}
	data := b.data
	b.data = nil
	b.cond.Broadcast()
	return data, nil
}

func (b *uploadBuffer) Close() error {
	b.mu.Lock()
	b.closed = true
	b.data = nil
	b.cond.Broadcast()
	b.mu.Unlock()
	return nil
}

func (c *Client) uploadPackets(ctx context.Context, conn *streamConn, lease *poolLease, session string) {
	maxSize := int(c.config.postSize.sample())
	buffer := newUploadBuffer(maxSize)
	defer buffer.Close()
	go func() {
		_, err := io.Copy(buffer, conn.wire)
		buffer.Close()
		if ctx.Err() == nil {
			conn.finish(err)
		}
	}()
	stop := context.AfterFunc(ctx, func() { buffer.Close() })
	defer stop()
	slots := make(chan struct{}, maxConcurrentPosts)
	var last time.Time
	var seq uint64
	for {
		interval := time.Duration(c.config.postInterval.sample()) * time.Millisecond
		if delay := time.Until(last.Add(interval)); delay > 0 {
			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
		select {
		case slots <- struct{}{}:
		case <-ctx.Done():
			return
		}
		data, err := buffer.take()
		if err != nil {
			return
		}
		req := c.config.newRequest(ctx, session, strconv.FormatUint(seq, 10), nil, data)
		seq++
		last = time.Now()
		go func() {
			defer func() { <-slots }()
			response, err := lease.RoundTrip(req)
			if err != nil {
				conn.finish(err)
				return
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				conn.finish(fmt.Errorf("xhttp: upload HTTP status: %s", response.Status))
				return
			}
			// Bound unexpected acknowledgement bodies rather than draining indefinitely.
			_, err = io.Copy(io.Discard, io.LimitReader(response.Body, 65536))
			if err != nil {
				conn.finish(err)
			}
		}()
	}
}
