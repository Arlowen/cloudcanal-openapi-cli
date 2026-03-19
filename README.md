# cloudcanal-openapi-cli

基于 Go 实现的 CloudCanal OpenAPI CLI，支持交互式使用，也支持单次命令执行。

当前已支持的 `datajob` 操作：

- `list`
- `show`
- `start`
- `stop`
- `delete`
- `replay`

## 快速开始

要求：

- Go 1.25+

构建并测试：

```bash
./scripts/all_build.sh
```

安装到命令行：

```bash
./scripts/install.sh
```

卸载：

```bash
./scripts/uninstall.sh
```

如果想看完整构建输出：

```bash
VERBOSE=1 ./scripts/all_build.sh
```

## 使用方式

交互模式：

```bash
cloudcanal
```

单次命令模式：

```bash
cloudcanal jobs list
cloudcanal jobs show 123
cloudcanal jobs replay 123 --auto-start
```

如果还没有执行安装脚本，也可以直接运行二进制：

```bash
./bin/cloudcanal jobs list
```

## 初始化配置

第一次启动会进入初始化向导，配置文件保存到：

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

说明：

- `apiBaseUrl` 必须是完整 URL，包含 `http://` 或 `https://`
- `secretKey` 不会在 `config show` 中明文展示

## 命令

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

## 开发

只编译：

```bash
make build
```

只测试：

```bash
make test
```
