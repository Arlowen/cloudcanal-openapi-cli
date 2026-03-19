# cloudcanal-openapi-cli

CloudCanal OpenAPI 的命令行工具，支持：

- 交互式命令行
- 单次命令执行
- `--output json` 机器可读输出
- zsh / bash TAB 补全

完整命令说明见 [docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)。

## 快速开始

安装：

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_install.sh | bash
```

说明：

- 安装脚本会从 GitHub Releases 下载预编译二进制
- 下载后会自动校验 `checksums.txt`
- 不需要本机安装 Go
- 会自动安装命令、PATH 和 zsh / bash 补全

安装完成后，先直接运行：

```bash
cloudcanal
```

首次启动会进入初始化向导，配置完成后就可以开始用。

## 常用用法

交互模式：

```bash
cloudcanal
```

单次命令：

```bash
cloudcanal jobs list
cloudcanal jobs show 123
cloudcanal datasources list --type MYSQL
cloudcanal workers list --cluster-id 2
```

JSON 输出：

```bash
cloudcanal jobs list --type SYNC --output json
```

## 配置

配置文件默认保存在：

```text
~/.cloudcanal/config.json
```

最小配置示例：

```json
{
  "apiBaseUrl": "https://cc.example.com",
  "accessKey": "your-ak",
  "secretKey": "your-sk",
  "language": "en"
}
```

如果你需要调整网络行为，也可以追加这些可选项：

```json
{
  "httpTimeoutSeconds": 15,
  "httpReadMaxRetries": 2,
  "httpReadRetryBackoffMillis": 300
}
```

## 文档入口

- 安装、初始化、命令参数、示例：[docs/cloudcanal-cli-usage.md](docs/cloudcanal-cli-usage.md)
- 机器可读输出：在查询命令后追加 `--output json`
- 补全脚本：`cloudcanal completion zsh` / `cloudcanal completion bash`

## 卸载

```bash
curl -fsSL https://raw.githubusercontent.com/Arlowen/cloudcanal-openapi-cli/main/scripts/bootstrap_uninstall.sh | bash
```

## 开发

要求：

- Go 1.25+

常用命令：

```bash
./scripts/all_build.sh
make build
make test
./scripts/install.sh
./scripts/uninstall.sh
```

发布：

- 推送 tag，例如 `v0.1.0`
- GitHub Actions 会自动构建并发布 release 资产
- Release 会同时生成 `checksums.txt`
