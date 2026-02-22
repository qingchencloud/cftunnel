# cftunnel — Cloudflare Tunnel 管理

快速管理 Cloudflare Tunnel 的 CLI 工具，3 条命令搞定内网穿透。

## 安装

```bash
# macOS / Linux
curl -fsSL https://github.com/qingchencloud/cftunnel/releases/latest/download/cftunnel_$(uname -s | tr A-Z a-z)_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o /usr/local/bin/cftunnel && chmod +x /usr/local/bin/cftunnel
```

## 快速上手

```bash
cftunnel init                                          # 交互式初始化
cftunnel add myapp 3000 --domain myapp.example.com     # 添加路由
cftunnel up                                            # 启动隧道
```

## 命令参考

| 命令 | 说明 |
|------|------|
| `init` | 交互式初始化（支持 `--token`/`--account`/`--name` 非交互模式） |
| `add <名称> <端口> --domain <域名>` | 添加路由（自动创建 CNAME） |
| `remove <名称>` | 删除路由（清理 DNS） |
| `list` | 列出所有路由 |
| `up` | 启动隧道 |
| `down` | 停止隧道 |
| `status` | 查看状态 |
| `destroy [--force]` | 删除隧道 + 所有 DNS 记录 |
| `reset [--force]` | 完全重置（删隧道 + 清本地配置） |
| `install` | 注册系统服务（开机自启） |
| `uninstall` | 卸载系统服务 |
| `logs [-f]` | 查看日志 |
| `version [--check]` | 版本信息 / 检查更新 |
| `update` | 自动更新到最新版 |

## CF API Token 权限

```
帐户 │ Cloudflare Tunnel │ 编辑
区域 │ DNS              │ 编辑
区域 │ 区域设置          │ 读取
```
