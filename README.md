# cloudcanal-openapi-cli

CloudCanal OpenAPI 的交互式 CLI，Go 实现。

当前 MVP 支持：

- 首次启动时录入并保存 `apiHost`、`ak`、`sk`
- 查看任务列表
- 查看任务详情
- 启动任务
- 停止任务
- 删除任务
- 重放任务

## 要求

- Go 1.25+

## 构建

```bash
go build -o bin/cloudcanal ./cmd/cloudcanal
```

或者：

```bash
make build
```

或者一键清理并测试、编译：

```bash
./scripts/all_build.sh
```

如果想看完整测试和构建输出：

```bash
VERBOSE=1 ./scripts/all_build.sh
```

安装到命令行环境：

```bash
./scripts/install.sh
```

卸载：

```bash
./scripts/uninstall.sh
```

## 运行

```bash
./bin/cloudcanal
```

也支持直接执行单条命令：

```bash
./bin/cloudcanal jobs list
./bin/cloudcanal jobs show 123
./bin/cloudcanal jobs replay 123 --auto-start
```

第一次启动如果不存在配置文件，会进入初始化向导：

```text
CloudCanal CLI initialization
Type exit at any prompt to cancel.
apiHost must be a full URL, for example: https://cc.example.com
apiHost:
ak:
sk:
```

配置文件保存到：

```text
~/.cloudcanal/config.json
```

配置格式：

```json
{
  "apiBaseUrl": "https://cc.example.com",
  "accessKey": "your-ak",
  "secretKey": "your-sk"
}
```

`apiBaseUrl` 必须是完整 URL，包含 `http://` 或 `https://`。

## 命令

进入 CLI 后可用命令：

```text
jobs list
jobs show <jobId>
jobs start <jobId>
jobs stop <jobId>
jobs delete <jobId>
jobs replay <jobId> [--auto-start] [--reset-to-created]
config show
config init
help
exit
quit
```

## 测试

```bash
go test ./...
```

或者：

```bash
make test
```
