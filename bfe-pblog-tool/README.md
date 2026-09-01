# bfePblogTool 使用说明书

## 使用

### bfePblogTool help 

```bash
%>./bfePblogTool -h                         

NAME:
   bfePblogTool

USAGE:
   bfePblogTool [global options] command [command options] [arguments...]

VERSION:
   v0.3.6

DESCRIPTION:
   https://github.com/bfenetworks/bfe-access-pb

COMMANDS:
   cat      举例：pbtool cat /tmp/pb_access3.log
   tail     示例: pbtool tail -n 20 [-f --interval 200] /path/to/pb_access3.log
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### bfePblogTool cat

`cat` 命令用于顺序输出 PB3 日志文件中的所有记录。可选 `-n` 参数显示行号，适合快速全量浏览或与其它命令（grep/awk）组合。

#### 用法

```bash
%> ./bfePblogTool cat -h
NAME:
   bfePblogTool cat - 举例：pbtool cat /tmp/pb_access3.log

USAGE:
   bfePblogTool cat [command options] [arguments...]

OPTIONS:
   -n  显示行号
```

#### 示例

```bash
# 全量输出日志
%> ./bfePblogTool cat /tmp/pb_access3.log

# 带行号输出
%> ./bfePblogTool cat -n /tmp/pb_access3.log
```

#### 注意事项

- 对超大文件（数 GB 级别）全量输出会产生大量 IO 与终端滚动，建议结合管道过滤或使用 `tail` 查看末尾
- 若需要实时新增记录，请使用 `tail -f` 而不是 `cat`

### bfePblogTool tail

`tail` 命令用于快速查看 PB3 日志文件末尾的若干条记录，并支持持续跟随新增内容（类似 `tail -f`），在超大文件 (>=10GB) 下依然保持常量级内存。

#### 用法

```bash
%> ./bfePblogTool tail -h
NAME:
   bfePblogTool tail - 示例: pbtool tail -n 20 [-f --interval 200] /path/to/pb_access3.log

USAGE:
   bfePblogTool tail [command options] [arguments...]

OPTIONS:
   -n value          最后 N 条 (默认10) (default: 0)
   -f                持续跟随输出新增日志
   --interval value  跟随模式下的轮询间隔(毫秒), 默认500 (default: 0)
```

#### 典型示例

```bash
# 查看最后 10 条（默认）
%> ./bfePblogTool tail /var/log/pb_access3.log

# 查看最后 50 条
%> ./bfePblogTool tail -n 50 /var/log/pb_access3.log

# 持续跟随新增记录
%> ./bfePblogTool tail -f /var/log/pb_access3.log
```

---

本项目采用 [Apache License 2.0](../LICENSE) 开源协议。
