### 结构

可用于入站和出站的 `transport` 字段。`splithttp` 是 `xhttp` 的别名。

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

### 字段

范围字段可使用整数或包含两端的范围，例如 `"100-1000"`。省略或范围两端均为零时使用默认值。

#### host

HTTP 主机名。客户端默认使用 TLS 服务器名称或服务器地址。设置后，服务端会验证该值。

#### path

HTTP 请求路径。默认使用 `/`。客户端与服务端必须一致。

#### mode

| 模式 | 上传 | 下载 |
| --- | --- | --- |
| `packet-up` | 多个 POST 请求 | 流式 GET |
| `stream-up` | 流式 POST | 流式 GET |
| `stream-one` | 流式 POST | 同一 POST 的响应 |

默认使用 `auto`。客户端通常使用 `packet-up`，与 REALITY 配合时使用 `stream-one`。
除非配置了指定模式，否则服务端接受所有模式。

`stream-up` 和 `stream-one` 需要 HTTP/2 或 HTTP/3。

#### http_version

可用值：`1.1`、`2`、`3`。

未启用 TLS 时，客户端默认使用 HTTP/1.1；启用 TLS 时默认使用 HTTP/2。TLS ALPN 仅包含 `http/1.1` 或 `h3` 时，会选择对应版本。
服务端默认接受 HTTP/1.1 和 HTTP/2；设置版本后仅接受该版本。

HTTP/3 需要 TLS 和包含 QUIC 支持的构建。REALITY 使用 HTTP/2。

#### headers

客户端使用额外请求标头，服务端使用额外响应标头。请使用 `host` 设置 HTTP 主机名。

#### x_padding_bytes

填充大小，上限为 65536。默认使用 `"100-1000"`。客户端生成的大小必须在服务端接受的范围内。

#### x_padding_obfs_mode

启用下列自定义填充设置。默认禁用；禁用时会在 `Referer` 中发送重复的 `X` 填充。

#### x_padding_method

`repeat-x` 重复字符 `X`；`tokenish` 生成随机字母数字填充。
`repeat-x` 的大小按字符数计算，`tokenish` 的大小按 Huffman 编码后的字节数计算。

默认使用 `repeat-x`。

#### x_padding_placement

| 值 | 位置 |
| --- | --- |
| `queryInHeader` | `x_padding_header` 中的 URL 查询字符串（默认） |
| `header` | `x_padding_header` 的值 |
| `cookie` | 名称为 `x_padding_key` 的 Cookie |
| `query` | 名称为 `x_padding_key` 的请求查询参数 |

#### x_padding_key

填充查询参数或 Cookie 名称。默认使用 `x_padding`。

#### x_padding_header

填充标头名称。默认使用 `X-Padding`。

#### session_id_placement / seq_placement

会话 ID 和数据包序号的位置。可用值：`path`、`query`、`header`、`cookie`。
默认使用 `path`。客户端与服务端必须一致。

#### session_id_key / seq_key

所选位置使用的名称。选择 `path` 时忽略。客户端与服务端必须一致。

标头默认使用 `X-Session` / `X-Seq`；查询参数与 Cookie 默认使用 `x_session` / `x_seq`。

#### session_id_table

客户端会话 ID 字符表。留空时使用 UUID v4。

可用字符表：`ALPHABET`、`alphabet`、`Alphabet`、`BASE36`、`base36`、`Base62`、`HEX`、`hex`、`number`。
自定义字符表可以包含互不重复的字母、数字以及 `-`、`.`、`_`、`~`。

#### session_id_length

使用 `session_id_table` 时必填。可使用 1–4096 之间的整数或范围；最短长度必须提供至少 31 位随机性。

#### sc_max_buffered_posts

服务端的数据包重排序窗口，范围为 1–1024。默认使用 `30`。

#### no_sse_header

仅服务端。下载响应中不发送 `Content-Type: text/event-stream`。默认禁用。

#### no_grpc_header

仅客户端。流式上传中不发送 `Content-Type: application/grpc`。默认禁用。

#### sc_max_each_post_bytes

数据包上传的最大大小，上限为 16 MiB。默认使用 `1000000`。
客户端为每个连接从范围中取值；服务端接受不超过范围上限的值。

#### sc_min_posts_interval_ms

仅客户端。数据包上传的最小间隔，单位为毫秒，上限为 60000。默认使用 `30`。

#### xmux

客户端 HTTP 连接复用设置，与协议多路复用相互独立。

| 字段 | 描述 |
| --- | --- |
| `max_connections` | 目标连接数，上限为 1024。与 `max_concurrency` 冲突。 |
| `max_concurrency` | 每个连接的最大并发流数量。 |
| `c_max_reuse_times` | 每个连接可分配的最大流数量。 |
| `h_max_request_times` | 每个连接的最大 HTTP 请求数。 |
| `h_max_reusable_secs` | 连接可复用时间，单位为秒。 |
| `h_keep_alive_period` | 保活间隔，单位为秒。`0` 使用 30 秒；`-1` 禁用。 |

计数字段可使用整数或范围。省略 `xmux` 或使用空对象时，使用结构示例中的默认值。
在非空对象中，未设置的计数字段不受限制。
