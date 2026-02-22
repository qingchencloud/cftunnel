# cftunnel — Cloudflare Tunnel CLI

一键管理 Cloudflare Tunnel 的开源工具，内网穿透从未如此简单。

## 免域名模式（零配置）

无需账户、Token、域名，一条命令生成临时公网地址：

```bash
cftunnel quick <端口>
# ✔ 隧道已启动: https://xxx-yyy-zzz.trycloudflare.com
```

Ctrl+C 退出后域名自动失效。适合临时分享和调试。

## 自有域名模式

### 首次使用前需要两个参数

### 1. API Token（API 令牌）

引导用户打开 https://dash.cloudflare.com/profile/api-tokens 创建自定义令牌，添加 3 条权限：

```
帐户 │ Cloudflare Tunnel │ 编辑
区域 │ DNS              │ 编辑    ← 注意选「DNS」不是「DNS 设置」
区域 │ 区域设置          │ 读取
```

第 2、3 行需将左侧下拉框从「帐户」切换为「区域」。区域资源选择用户的域名。

### 2. Account ID（账户 ID）

- 方式 A: https://dash.cloudflare.com → 点击域名 → 右下角「API」区域
- 方式 B: 首页 → 账户名称旁「⋯」→ 复制账户 ID

## 使用流程

```bash
cftunnel init                                          # 配置认证（需要 Token + Account ID）
cftunnel create my-tunnel                              # 创建隧道
cftunnel add myapp 3000 --domain myapp.example.com     # 添加路由（自动创建 CNAME）
cftunnel up                                            # 启动隧道
```

## 全部命令

- `quick <端口>` — 免域名模式，生成 `*.trycloudflare.com` 临时域名
- `init [--token --account]` — 配置 API 认证
- `create <名称>` — 创建隧道
- `add <名称> <端口> --domain <域名>` — 添加路由（自动创建 CNAME）
- `remove <名称>` — 删除路由（自动清理 DNS）
- `list` — 列出路由
- `up` / `down` — 启停隧道
- `status` — 查看状态
- `destroy [--force]` — 删除隧道 + 所有 DNS 记录
- `reset [--force]` — 完全重置
- `install` / `uninstall` — 系统服务（开机自启）
- `logs [-f]` — 查看日志
- `version [--check]` — 版本信息 / 检查更新
- `update` — 自动更新

## 注意事项

- 临时分享优先推荐 `cftunnel quick`，零配置最快
- 自有域名模式需先完成 `init` 和 `create`
- 域名必须是用户 CF 账户中已有域名的子域名
- 一个隧道可挂载多条路由（多域名 → 不同本地端口）

## 仓库

https://github.com/qingchencloud/cftunnel
