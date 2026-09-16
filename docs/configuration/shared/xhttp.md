### Structure

Available in inbound and outbound `transport` fields. `splithttp` is an alias for `xhttp`.

```json
{
  "type": "xhttp",
  "host": "",
  "path": "/",
  "mode": "auto",
  "http_version": "",
  "headers": {},
  "x_padding_bytes": "100-1000",
  "x_padding_obfs_mode": false,
  "x_padding_method": "repeat-x",
  "x_padding_placement": "queryInHeader",
  "x_padding_key": "x_padding",
  "x_padding_header": "X-Padding",
  "session_id_placement": "path",
  "session_id_key": "",
  "session_id_table": "Base62",
  "session_id_length": 16,
  "seq_placement": "path",
  "seq_key": "",
  "sc_max_buffered_posts": 30,
  "no_sse_header": false,
  "no_grpc_header": false,
  "sc_max_each_post_bytes": 1000000,
  "sc_min_posts_interval_ms": 30,
  "xmux": {
    "max_connections": 3,
    "h_max_request_times": "600-900",
    "h_max_reusable_secs": "1800-3000"
  }
}
```

### Fields

Range fields accept an integer or an inclusive range such as `"100-1000"`. Omitted or all-zero ranges use the default.

#### host

HTTP host. The client defaults to the TLS server name or server address. The server verifies it if set.

#### path

HTTP request path. `/` is used by default. Must match on both sides.

#### mode

| Mode | Upload | Download |
| --- | --- | --- |
| `packet-up` | Multiple POST requests | Streaming GET |
| `stream-up` | Streaming POST | Streaming GET |
| `stream-one` | Streaming POST | Response to the same POST |

`auto` is used by default. The client uses `packet-up`, or `stream-one` with REALITY.
The server accepts all modes unless a specific mode is set.

`stream-up` and `stream-one` require HTTP/2 or HTTP/3.

#### http_version

Available values: `1.1`, `2`, `3`.

The client defaults to HTTP/1.1 without TLS and HTTP/2 with TLS. A single TLS ALPN value of `http/1.1` or `h3` selects that version.
The server accepts HTTP/1.1 and HTTP/2 by default; setting a version restricts it to that version.

HTTP/3 requires TLS and a build with QUIC support. REALITY uses HTTP/2.

#### headers

Extra request headers on the client, response headers on the server. Use `host` to set the HTTP host.

#### x_padding_bytes

Padding size, up to 65536. `"100-1000"` is used by default. Must fit the server's accepted range.

#### x_padding_obfs_mode

Enable custom padding settings below. Disabled by default; repeated-X padding is sent in `Referer` instead.

#### x_padding_method

`repeat-x` repeats `X`; `tokenish` generates random alphanumeric padding.
Size is measured in characters for `repeat-x`, Huffman-encoded bytes for `tokenish`.

`repeat-x` is used by default.

#### x_padding_placement

| Value | Placement |
| --- | --- |
| `queryInHeader` | URL query in `x_padding_header` (default) |
| `header` | Value in `x_padding_header` |
| `cookie` | Cookie named `x_padding_key` |
| `query` | Request query parameter named `x_padding_key` |

#### x_padding_key

Padding query parameter or cookie name. `x_padding` is used by default.

#### x_padding_header

Padding header name. `X-Padding` is used by default.

#### session_id_placement / seq_placement

Session ID and packet sequence placement. Available values: `path`, `query`, `header`, `cookie`.
`path` is used by default. Must match on both sides.

#### session_id_key / seq_key

Names used for the selected placement. Ignored for `path`. Must match on both sides.

Defaults: `X-Session` / `X-Seq` for headers, `x_session` / `x_seq` for query parameters and cookies.

#### session_id_table

Client session ID alphabet. UUID v4 is used when empty.

Available tables: `ALPHABET`, `alphabet`, `Alphabet`, `BASE36`, `base36`, `Base62`, `HEX`, `hex`, `number`.
A custom table may contain unique letters, digits, `-`, `.`, `_`, and `~`.

#### session_id_length

Required with `session_id_table`. Integer or range within 1–4096; the shortest length must provide at least 31 bits of randomness.

#### sc_max_buffered_posts

Server packet reordering window, within 1–1024. `30` is used by default.

#### no_sse_header

Server only. Omit `Content-Type: text/event-stream` on downloads. Disabled by default.

#### no_grpc_header

Client only. Omit `Content-Type: application/grpc` on streaming uploads. Disabled by default.

#### sc_max_each_post_bytes

Maximum packet upload size, up to 16 MiB. `1000000` is used by default.
The client samples the range per connection; the server accepts up to its upper bound.

#### sc_min_posts_interval_ms

Client only. Minimum interval between packet uploads, in milliseconds, up to 60000. `30` is used by default.

#### xmux

Client HTTP connection reuse settings, independent of protocol multiplexing.

| Field | Description |
| --- | --- |
| `max_connections` | Target connection count, up to 1024. Conflicts with `max_concurrency`. |
| `max_concurrency` | Maximum concurrent streams per connection. |
| `c_max_reuse_times` | Maximum stream assignments per connection. |
| `h_max_request_times` | Maximum HTTP requests per connection. |
| `h_max_reusable_secs` | Connection reuse lifetime in seconds. |
| `h_keep_alive_period` | Keepalive interval in seconds. `0` uses 30 seconds; `-1` disables it. |

Counters accept integers or ranges. An omitted or empty `xmux` uses the values shown above.
In a nonempty object, unspecified counters are unlimited.
