# Beaker v3 Configuration Migration

目标版本：Go 1.27.1。

## 结论

本次变更不需要修改线上 MySQL 表结构，也不需要新增或删除服务器 TOML 字段。

需要注意的变化只有两类：

- Redis 缓存可以清理，v3 会按 `REDIS_PREFIX` 重新生成缓存。
- CLI 本地配置会新增 `SessionToken`，旧 CLI 不再作为兼容目标。

## 服务器 TOML

`config.toml`：无字段变化。

`admin.toml`：无字段变化。

线上服务器不需要因为本次代码改动手动增加 TOML 字段。仍然需要确认以下原有字段正确：

- `admin.toml` 的 `[authinfo].EXPIRE_TIME`：session token 的有效期仍使用这个值。
- `admin.toml` 的 `[authinfo].ServerKeyDir`：目录需要可写，服务启动时会生成或读取 `pub.pem` 和 `pri.pem`。
- `config.toml` 和 `admin.toml` 的 `[redis].REDIS_PREFIX`：必须是当前 Beaker 实例专用前缀，推荐继续使用 `beaker_` 或每个环境独立前缀。
- `config.toml` 和 `admin.toml` 的 `[database]`：MySQL 地址、账号、密码、库名和连接池配置保持不变。

## MySQL

MySQL 库结构不能动，本次变更没有新增 migration，也没有修改 `build/beaker.sql`。

发布前建议只做常规备份，不需要执行任何 DDL：

```sql
-- 不需要执行 ALTER TABLE。
-- 只建议按线上流程做发布前备份。
```

## Redis

线上 Redis 可以清理缓存内容。

v3 之前的风险：清缓存使用 `FLUSHALL`，会清掉同一个 Redis 实例里的所有业务 key。

v3 行为：清缓存只扫描并删除：

```text
REDIS_PREFIX + "*"
```

例如：

```toml
[redis]
REDIS_PREFIX = "beaker_"
```

只会删除 `beaker_*`。

手动清理示例：

```bash
redis-cli --scan --pattern 'beaker_*' | xargs -r redis-cli del
```

如果线上 Redis 是共享实例，发布前必须确认 `REDIS_PREFIX` 不会和其他应用冲突。

## CLI 本地配置

v3 CLI 登录成功后，本地 `~/.beaker` 中每个站点会新增：

```yaml
SessionToken: "..."
```

用户需要重新登录：

```bash
beaker login -u <admin-name> -p <admin-password>
```

或者重新添加站点：

```bash
beaker config addw <admin-url> -a <alias> -u <admin-name> -p <admin-password> -d
```

旧 CLI 不再保留兼容。没有 `SessionToken` 的管理请求会返回：

```text
Error Code: BEAKER-401
Description: login has expired or the session token is invalid. Run `beaker login` again.
Detail: Need Login
Request ID: <server-request-id>
```

## 线上发布步骤

1. 在发布环境安装或切换到 Go 1.27.1。
2. 拉取 `v3dev` 分支代码。
3. 在有网络环境执行 `go mod tidy` 并确认 `go.sum` 完整。
4. 执行 `go test ./...`、`go vet ./...`、`govulncheck ./...`。
5. 构建新的 server、admin、cli 二进制。
6. 确认服务器 TOML 路径环境变量仍正确：`BEAKERPATH` 和 `BEAKERADMINPATH`。
7. 发布 server 和 admin。
8. 清理 Redis 中 Beaker 前缀缓存，或等待缓存自然过期。
9. 用户更新 CLI 后重新执行 `beaker login`。

## 回滚注意事项

- MySQL 未改结构，因此数据库无需回滚。
- Redis 缓存可删除后由旧版本重新生成。
- 如果回滚到旧 CLI/旧 admin，旧 CLI 不理解 `SessionToken`，但本地 YAML 中多出的字段不会影响读取；仍建议重新登录一次。
