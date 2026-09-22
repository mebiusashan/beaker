# Beaker v3 Technical Design

目标版本：Go 1.27.1。

## 设计目标

- 修复已知 Dependabot 安全告警，并降低后续升级成本。
- 服务端所有管理接口异常都返回稳定 `errorCode`，CLI 直接展示错误码和详细解释。
- 不修改 MySQL 表结构。
- Redis 缓存允许清理并自动重建，但不能误清共享 Redis 中其他业务的 key。
- 新 CLI 不保留旧 CLI 兼容；v3 管理接口要求 session token。

## 错误协议

服务端失败响应保持原有 `code/msg/data` 结构，同时新增：

```json
{
  "code": 1,
  "errorCode": "BEAKER-500-DB",
  "requestId": "f1c2...",
  "msg": "dial tcp 127.0.0.1:3306: connect: connection refused"
}
```

CLI 失败时直接输出：

```text
Error Code: BEAKER-500-DB
Description: database operation failed. Check MySQL address, credentials, network, connection limits, and whether MySQL is out of disk or memory.
Detail: dial tcp 127.0.0.1:3306: connect: connection refused
Request ID: f1c2...
```

## 错误码表

| Error Code | CLI 说明 |
| --- | --- |
| `BEAKER-400` | 请求格式不合法，检查 CLI 参数或生成的请求数据。 |
| `BEAKER-400-DECODE` | 请求加密或解密失败，通常需要重新 `beaker login`。 |
| `BEAKER-400-UPLOAD` | 图片上传被拒绝，检查文件名、后缀、MIME 和 5 MB 限制。 |
| `BEAKER-401` | 登录过期或 session token 无效，重新执行 `beaker login`。 |
| `BEAKER-403` | 当前用户无权限执行该操作。 |
| `BEAKER-404` | 资源不存在，检查 id、alias 或 URL。 |
| `BEAKER-500-DB` | 数据库操作失败，检查 MySQL 地址、账号、网络、连接数、磁盘和内存。 |
| `BEAKER-500-CACHE` | Redis 操作失败，检查 Redis 地址、网络、内存策略和可用内存。 |
| `BEAKER-500` | 服务端内部错误，用 request id 定位日志。 |
| `BEAKER-CLI-NETWORK` | CLI 无法连接服务端，检查 URL、DNS、防火墙和网络。 |
| `BEAKER-CLI-HTTP` | 服务端返回 HTTP 错误，检查 admin 服务和反向代理。 |
| `BEAKER-CLI-RESPONSE` | 服务端响应不是 Beaker JSON，通常是 URL 指错或反代返回了 HTML 错误页。 |

## 认证变化

v3 管理端登录成功后会返回：

- `data`: 仍是用于当前加密协议的服务端 key。
- `sessionToken`: 新增，CLI 后续请求会通过 `X-Beaker-Session` header 发送。

管理接口不再接受旧 CLI 只依赖全局登录密钥的请求。线上升级后，用户需要重新执行 `beaker login` 或重新添加站点配置。

## 安全改动

- HTTP Server 增加 `ReadHeaderTimeout`、`ReadTimeout`、`WriteTimeout`、`IdleTimeout` 和 `MaxHeaderBytes`，降低慢请求和异常连接风险。
- Gin 增加 request id、请求体大小限制和 panic recovery 中间件。
- 图片上传只接受固定后缀和匹配 MIME，最大 5 MB，写入使用临时文件后原子 rename。
- Markdown 渲染跳过原始 HTML 和 style，并启用 safe-link 过滤，降低存储型 XSS 风险。
- 3DES 解密遇到畸形密文时返回错误，不再 panic。
- session key 使用 `crypto/rand` 生成。

## 性能改动

- Redis 改为连接池，避免每次缓存读写都创建 TCP 连接。
- Jet 模板集合缓存到 `ViewRender`，避免每次渲染重新扫描模板目录。
- Redis 清理改为按 `REDIS_PREFIX` 扫描删除，避免 `FLUSHALL` 阻塞并误清其他业务。
- Tweet 空列表分页不再访问 page 0，避免无意义查询和 offset 下溢。

## 数据兼容

- MySQL：不修改表结构，不需要迁移 SQL。
- Redis：可以全部清理。v3 代码只会删除 `REDIS_PREFIX` 命名空间内的缓存 key。
- TOML：服务器 TOML 字段未变化。
- CLI 本地配置：新增 `SessionToken`，由 `beaker login` 自动写入。

## 依赖说明

已在可访问 Go module proxy 的环境中执行依赖升级和 `go mod tidy`，当前锁定版本为：

- `github.com/gin-gonic/gin v1.12.0`
- `golang.org/x/net v0.59.0`
- `golang.org/x/crypto v0.57.0`
- `golang.org/x/sys v0.48.0`
- `golang.org/x/text v0.42.0`
- `google.golang.org/protobuf v1.36.12`

已验证 `GOCACHE=/tmp/beaker-gocache go test ./...` 通过。发布前仍建议在 CI 中继续执行 `go test -race ./...`、`go vet ./...` 和 `govulncheck ./...`。
