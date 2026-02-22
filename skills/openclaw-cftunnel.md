# cftunnel — Cloudflare Tunnel CLI

一键管理 Cloudflare Tunnel 的开源工具。

## 快速上手

```bash
cftunnel init                                          # 初始化
cftunnel add myapp 3000 --domain myapp.example.com     # 添加路由
cftunnel up                                            # 启动
```

## 全部命令

- `init [--token --account --name]` — 初始化隧道
- `add <名称> <端口> --domain <域名>` — 添加路由
- `remove <名称>` — 删除路由
- `list` — 列出路由
- `up` / `down` — 启停隧道
- `status` — 查看状态
- `destroy [--force]` — 删除隧道
- `reset [--force]` — 完全重置
- `install` / `uninstall` — 系统服务
- `logs [-f]` — 查看日志
- `version [--check]` — 版本信息
- `update` — 自动更新

## 仓库

https://github.com/qingchencloud/cftunnel
