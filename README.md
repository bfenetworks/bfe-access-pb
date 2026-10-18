# bfe-access-pb

[![License](https://img.shields.io/github/license/bfenetworks/bfe-access-pb)](./LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/bfenetworks/bfe-access-pb.svg)](https://pkg.go.dev/github.com/bfenetworks/bfe-access-pb)

BFE 访问日志（Access Log）的 Protocol Buffers 定义与二进制日志读写库（b2log），以独立 Go Module 形式发布，供 [BFE](https://github.com/bfenetworks/bfe) 及其他下游仓库直接依赖。

## 功能特性

- **Protobuf Schema**（`bfe_access_pb/bfe_access.proto`）：BFE 访问日志的字段定义，覆盖请求、响应、耗时、AI 观测（路由、计费、缓存等）等维度。
- **预生成 Go 代码**（`bfe_access_pb/bfe_access.pb.go`）：生成代码随仓库一起发布，下游直接引用即可，无需本地安装 `protoc`，降低编译依赖复杂度。
- **b2log 二进制日志库**（`b2log/`）：将 Protobuf 序列化数据封装为带魔数与头部的二进制记录，提供读写 API。
- **语义化版本管理**：通过 tag 发布，下游以 `go get` 指定版本升级字段定义。

## 仓库结构

| 目录/文件 | 说明 |
|-----------|------|
| `bfe_access_pb/` | 访问日志 Protobuf Schema（`bfe_access.proto`）及生成的 Go 代码（`bfe_access.pb.go`） |
| `b2log/` | 二进制日志记录读写库 |
| `docs/protobuf.md` | 访问日志字段完整清单与字段编号规则 |
| `docs/README.md` | b2log 二进制格式与使用示例 |
| `build.sh` | 由 `.proto` 重新生成 `.pb.go` 的脚本 |
| `CHANGELOG.md` | 版本变更记录（Keep a Changelog 格式） |

## 快速开始

```bash
go get github.com/bfenetworks/bfe-access-pb@v0.3.16
```

```go
import (
    "google.golang.org/protobuf/proto"

    "github.com/bfenetworks/bfe-access-pb/bfe_access_pb"
)

log := &bfe_access_pb.RequestLog{ /* ... */ }
buf, err := proto.Marshal(log)
// ...
```

b2log 记录的读写（`HeaderWrite` / `BuffParse`）用法见 [docs/README.md](./docs/README.md)。

## 重新生成 Protobuf 代码

修改 `bfe_access_pb/bfe_access.proto` 后执行：

```bash
sh build.sh
```

脚本使用 `protoc`（默认 `/opt/protoc`，可用 `PROTOC` 环境变量覆盖）安装 `protoc-gen-go@v1.35.0` 并重新生成 `bfe_access.pb.go`。**不要手工编辑生成文件**。

## 版本与升级

- 版本号遵循[语义化版本](https://semver.org/lang/zh-CN/)，变更记录见 [CHANGELOG.md](./CHANGELOG.md)。
- 下游仓库通过升级依赖版本号获取新的 `.go` 文件：

```bash
go get github.com/bfenetworks/bfe-access-pb@v0.3.16
```

- 已发布 tag 中的字段编号保持稳定的，详见 [docs/protobuf.md](./docs/protobuf.md) 的编号规则。

## 参与贡献

- 新增访问日志字段：编辑 `bfe_access.proto`，按 [docs/protobuf.md](./docs/protobuf.md) 的语义区间选择字段编号并更新文档，然后运行 `sh build.sh` 重新生成。
- 提交前请运行测试：

```bash
go test ./...
```

## License

本项目基于 [Apache License 2.0](./LICENSE) 开源。
