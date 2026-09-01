# BFE Access Log Protobuf Schema & Tools

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

BFE 访问日志的 Protocol Buffers 定义、Go 语言封装库，以及日志处理工具。

## 项目简介

本仓库是 [BFE (Beyond Front End)](https://github.com/bfenetworks/bfe) 七层负载均衡器的访问日志组件，提供：

- **Protocol Buffers 定义**：BFE 访问日志的结构化 Schema (`bfe_access.proto`)
- **Go 语言封装**：日志的二进制读写库 (`b2log`)
- **日志处理工具**：`bfePblogTool` 命令行工具，用于查看和分析 PB 日志文件

其他 BFE 相关仓库（如 `bfe` 主仓库）可直接依赖本模块，无需自行处理 Protobuf 生成和依赖管理。

## 目录结构

```
bfe-access-pb/
├── bfe_access_pb/          # Protobuf Schema 及生成的 Go 代码
│   ├── bfe_access.proto    # 访问日志字段定义
│   └── bfe_access.pb.go    # protoc-gen-go 生成的 Go 代码（勿手动编辑）
├── b2log/                  # 二进制日志记录读写库
│   ├── b2log.go            # 核心数据结构
│   ├── b2log_read.go       # 日志读取
│   └── b2log_write.go      # 日志写入
├── bfe-pblog-tool/         # 日志处理命令行工具
│   ├── cmd/                # CLI 入口（cat / tail 命令）
│   └── bfe_reader/         # 日志读取核心逻辑
├── docs/                   # 文档
│   ├── protobuf.md         # 字段详细说明
│   └── README.md           # b2log 二进制格式说明
├── Makefile                # 构建、测试、发布脚本
└── VERSION                 # 版本号
```

## 快速开始

### 环境要求

- Go 1.22+
- Protocol Buffers 编译器（protoc，用于代码生成）

### 安装

```bash
# 作为依赖引入
go get github.com/bfenetworks/bfe-access-pb
```

### 构建

```bash
# 构建日志处理工具
make build

# 运行测试
make test

# 交叉编译多平台发布包
make release
```

构建产物位于：
- 单平台：`output/bin/bfePblogTool`
- 多平台发布：`dist/bfePblogTool_vX.Y.Z_<os>_<arch>.tar.gz`

## 日志处理工具 (bfePblogTool)

`bfePblogTool` 是用于查看 BFE PB 格式日志的命令行工具，支持类似 `cat` 和 `tail` 的操作。

### 使用方法

```bash
# 查看帮助
./bfePblogTool -h

# 输出整个日志文件
./bfePblogTool cat <日志文件路径>

# 带行号输出
./bfePblogTool cat -n <日志文件路径>

# 查看最后 N 条记录
./bfePblogTool tail -n <数量> <日志文件路径>

# 持续监听新增日志（类似 tail -f）
./bfePblogTool tail -f <日志文件路径>
```

### 示例

```bash
# 全量输出
./bfePblogTool cat /var/log/bfe/pb_access3.log

# 查看最后 50 条
./bfePblogTool tail -n 50 /var/log/bfe/pb_access3.log

# 实时跟随
./bfePblogTool tail -f /var/log/bfe/pb_access3.log
```

### 特性

- 流式处理大文件（支持数 GB 级别）
- 常量内存占用（即使处理超大文件）
- 支持文件轮转检测（inode 变化时自动重新打开）
- 支持 Ctrl+C 优雅退出

详细说明请参考 [bfe-pblog-tool/README.md](bfe-pblog-tool/README.md)

## 开发指南

### 修改 Protobuf Schema

如需添加或修改访问日志字段：

1. 编辑 `bfe_access_pb/bfe_access.proto`
2. 运行 `make proto` 重新生成 Go 代码
3. 更新 `docs/protobuf.md` 文档
4. 运行 `make test` 验证

### 版本管理

- 版本号记录在 `VERSION` 文件中
- 修改 Schema 后需更新版本号并打 tag
- 其他仓库通过 `replace` 指令在开发阶段使用本地模块

## 文档

- [Protocol Buffers 字段说明](docs/protobuf.md)
- [b2log 二进制格式](docs/README.md)
- [bfePblogTool 使用手册](bfe-pblog-tool/README.md)

## 参与贡献

欢迎提交 Issue 和 Pull Request。贡献前请确保：

- 代码符合项目风格规范
- 新功能附带测试用例
- 修改 Protobuf Schema 后运行 `make proto` 并更新文档

详细贡献指南请参考 [BFE 主仓库 CONTRIBUTING.md](https://github.com/bfenetworks/bfe/blob/develop/CONTRIBUTING.md)

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 开源协议。

## 相关链接

- [BFE 主仓库](https://github.com/bfenetworks/bfe)
- [BFE 官方文档](https://www.bfe-networks.net)
- [Protocol Buffers](https://protobuf.dev)
