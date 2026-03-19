# CLI 对齐修复报告

## 背景

用户在 `workers list` 的中文输出中观察到列未对齐。为避免只修单点，本次基于 mock 数据对所有主要命令做了一轮巡检。

## Mock 验证环境

- 使用 `test/repl/alignment_test.go` 中的 fake runtime 组装中英混排、长 ID、IP、浮点数、中文描述等数据。
- 使用 `test/util/table_test.go` 直接校验底层表格格式化逻辑。
- 通过 `go test ./...` 和 `./scripts/all_build.sh` 统一验证。

## 巡检范围

已检查：

- `jobs list`
- `jobs show`
- `jobs schema`
- `datasources list`
- `datasources show`
- `clusters list`
- `workers list`
- `consolejobs show`
- `job-config specs`

不属于本次表格对齐问题范围：

- `help`
- `config show`
- `lang show`
- `jobs start|stop|delete|replay`
- `workers start|stop`
- `clear`

## 发现

1. 所有表格命令共用 `internal/util/table.go`，原实现按 `len(string)` 计算宽度，中文和宽字符会导致列边界漂移。
2. 中文模式下仍有少量表头使用英文硬编码，例如 `Cloud`、`Cluster`、`Private IP`、`Job Type`，整体观感不一致。
3. 详情类命令虽然不是表格，但 `label: value` 没有统一冒号列，长短标签混排时可读性一般。

## 修复

1. 表格宽度计算切换到终端显示宽度，统一按可视列宽补空格。
2. 列表命令的相关表头统一走 `label()`，支持中英文切换。
3. 详情类输出的标签列统一补齐，冒号位置对齐。

## 结果

- 列表命令在中文、英文和中英混排数据下保持列对齐。
- 详情类命令的字段标签列保持统一。
- 对齐相关问题新增了 mock 回归测试，后续改动更容易发现回归。
