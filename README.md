# cftunnel

[![GitHub release](https://img.shields.io/github/v/release/qingchencloud/cftunnel)](https://github.com/qingchencloud/cftunnel/releases)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen?style=flat&logo=go)](https://goreportcard.com/report/github.com/qingchencloud/cftunnel)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Cloudflare Tunnel 一键管理 CLI** — 让本地项目秒变公网可访问。

[为什么选 cftunnel？](#why) · [架构原理](#architecture) · [安装](#install) · [快速上手](#quickstart) · [命令参考](#commands) · [AI 助手集成](#ai) · [交流](#contact)

关联项目：[cftunnel-app 桌面客户端](https://github.com/qingchencloud/cftunnel-app)（[下载](https://github.com/qingchencloud/cftunnel-app/releases)） · [ClawApp](https://github.com/qingchencloud/clawapp) · [OpenClaw 中文翻译](https://github.com/1186258278/OpenClawChineseTranslation)

> 用 AI 写了个前端页面想给客户看？本地跑着 API 想让远程同事调试？开发环境需要接收 Webhook？
>
> 一条命令，你的 `localhost` 就有了公网域名。

**免域名模式**（零配置，临时分享）：

```bash
cftunnel quick 3000
# ✔ 隧道已启动: https://xxx-yyy-zzz.trycloudflare.com
```

<p align="center">
  <img src="docs/images/terminal-demo.gif" alt="cftunnel quick 演示" width="720">
</p>

**自有域名模式**（稳定持久）：4 条命令搞定 `init` → `create` → `add` → `up`

cftunnel 把 Cloudflare Tunnel 的繁琐流程封装成极简 CLI，**免费、安全、无需公网 IP**。

<h2 id="why">为什么选 cftunnel？</h2>

<p align="center">
  <img src="docs/images/compare-chart.gif" alt="功能对比" width="720">
</p>

| 对比项 | 原生 cloudflared | cftunnel |
|--------|-----------------|----------|
| 免域名穿透 | `cloudflared tunnel --url` 手动操作 | `cftunnel quick <端口>` 一条命令 |
| 创建隧道 | 登录浏览器 + 手动配置 | `cftunnel create my-tunnel` |
| DNS 记录 | 手动去 Dashboard 创建 CNAME | `cftunnel add` 自动创建 |
| 多路由管理 | 手动编辑 YAML 配置 | `cftunnel add/remove/list` |
| 开机自启 | 手动写 systemd/launchd | `cftunnel install` |
| 清理资源 | 手动删隧道 + 删 DNS + 删配置 | `cftunnel destroy` 一键清理 |
| AI 集成 | 无 | 内置 Skills，AI 助手直接管理 |
| GUI 管理 | 无 | [桌面客户端](https://github.com/qingchencloud/cftunnel-app)，可视化操作 |

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="architecture">架构原理</h2>

<p align="center">
  <img src="docs/images/architecture.gif" alt="cftunnel 架构原理" width="720">
</p>

cftunnel 的工作流程：

1. **cftunnel CLI** 作为管理器，负责配置管理和进程编排
2. 自动下载并启动 **cloudflared** 守护进程
3. cloudflared 与 **Cloudflare Edge** 建立加密隧道
4. 互联网用户通过 Cloudflare 全球边缘节点访问你的本地服务

整个过程无需公网 IP，流量经过 Cloudflare 全球 CDN 加速，自带 HTTPS。

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="features">特性</h2>

- **免域名模式** — `cftunnel quick <端口>`，零配置生成 `*.trycloudflare.com` 临时公网地址
- **极简操作** — `init` → `create` → `add` → `up`，4 步搞定自有域名穿透
- **自动 DNS** — 添加路由时自动创建 CNAME 记录，删除时自动清理
- **进程托管** — 自动下载 cloudflared，支持 macOS launchd / Linux systemd / Windows Service 开机自启
- **自动更新** — 内置版本检查和一键自更新
- **AI 友好** — 内置 Claude Code / OpenClaw Skills，AI 助手可直接管理隧道
- **跨平台** — 支持 macOS (Intel/Apple Silicon) + Linux (amd64/arm64) + Windows (amd64/arm64)
- **桌面客户端** — 不想敲命令？[cftunnel-app](https://github.com/qingchencloud/cftunnel-app) 提供可视化 GUI，仪表盘一键启停

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="install">安装</h2>

### 一键安装（推荐）

**macOS / Linux：**

```bash
curl -fsSL https://raw.githubusercontent.com/qingchencloud/cftunnel/main/install.sh | bash
```

**Windows（PowerShell）：**

```powershell
irm https://raw.githubusercontent.com/qingchencloud/cftunnel/main/install.ps1 | iex
```

### 手动下载

从 [Releases](https://github.com/qingchencloud/cftunnel/releases) 下载对应平台的二进制文件：

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/qingchencloud/cftunnel/releases/latest/download/cftunnel_darwin_arm64.tar.gz | tar xz -C /usr/local/bin/

# macOS Intel
curl -fsSL https://github.com/qingchencloud/cftunnel/releases/latest/download/cftunnel_darwin_amd64.tar.gz | tar xz -C /usr/local/bin/

# Linux amd64
curl -fsSL https://github.com/qingchencloud/cftunnel/releases/latest/download/cftunnel_linux_amd64.tar.gz | tar xz -C /usr/local/bin/

# Linux arm64
curl -fsSL https://github.com/qingchencloud/cftunnel/releases/latest/download/cftunnel_linux_arm64.tar.gz | tar xz -C /usr/local/bin/
```

**Windows：** 从 [Releases](https://github.com/qingchencloud/cftunnel/releases) 下载 `cftunnel_windows_amd64.zip`，解压后将 `cftunnel.exe` 放到 PATH 目录中。

### 从源码构建

```bash
git clone https://github.com/qingchencloud/cftunnel.git
cd cftunnel
make build
```

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="quickstart">快速上手</h2>

### 方式一：免域名模式（零配置）

无需账户、Token、域名，装好就能用：

```bash
cftunnel quick 3000
# ✔ 隧道已启动: https://xxx-yyy-zzz.trycloudflare.com
```

> 适合临时分享和调试，Ctrl+C 退出后域名自动失效。需要固定域名请用方式二。

### 方式二：自有域名模式

<p align="center">
  <img src="docs/images/custom-domain.gif" alt="自有域名模式演示" width="720">
</p>

#### 1. 准备 Cloudflare API Token

> 前提：你需要一个 Cloudflare 账户和至少一个已添加的域名。

**创建 API 令牌：**

1. 打开 [API 令牌页面](https://dash.cloudflare.com/profile/api-tokens)
2. 点击「创建令牌」→「创建自定义令牌」→「开始使用」
3. 令牌名称随意填写（如 `cftunnel`）
4. 添加 3 条权限（点「+ 添加更多」逐条添加）：

```
┌──────────────────────────────────────────────────┐
│ 第 1 行: 帐户 │ Cloudflare Tunnel │ 编辑       │
│ 第 2 行: 区域 │ DNS              │ 编辑       │
│ 第 3 行: 区域 │ 区域设置          │ 读取       │
└──────────────────────────────────────────────────┘
```

> **注意**: 第 2、3 行需先将左侧下拉框从「帐户」切换为「区域」。DNS 权限请选择「DNS」而非「DNS 设置」，两者不同。

5. 区域资源 → 包括 → 特定区域 → 选择你的域名
6. 点击「继续以显示摘要」→「创建令牌」→ **立即复制令牌**（只显示一次）

**获取账户 ID（任选其一）：**

- **方式 A**: [Cloudflare 首页](https://dash.cloudflare.com) → 点击域名 → 页面右下角「API」区域 → 复制「账户 ID」
- **方式 B**: 首页 → 账户名称旁「⋯」→ 复制账户 ID

#### 2. 初始化认证

```bash
# 交互式（推荐）
cftunnel init

# 非交互式
cftunnel init --token <your-token> --account <account-id>
```

#### 3. 创建隧道

```bash
cftunnel create my-tunnel
```

#### 4. 添加路由

```bash
# 将 app.example.com 指向本地 3000 端口
cftunnel add myapp 3000 --domain app.example.com
```

#### 5. 启动隧道

```bash
cftunnel up
```

搞定！现在可以通过 `app.example.com` 访问你本地的 3000 端口服务了。

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="commands">命令参考</h2>

### 免域名模式

| 命令 | 说明 |
|------|------|
| `cftunnel quick <端口>` | 零配置启动，生成 `*.trycloudflare.com` 临时域名 |

### 配置管理

| 命令 | 说明 |
|------|------|
| `cftunnel init` | 配置认证信息（支持 `--token`/`--account`） |
| `cftunnel create <名称>` | 创建隧道 |
| `cftunnel add <名称> <端口> --domain <域名>` | 添加路由（自动创建 CNAME） |
| `cftunnel remove <名称>` | 删除路由（自动清理 DNS） |
| `cftunnel list` | 列出所有路由 |

### 运行管理

| 命令 | 说明 |
|------|------|
| `cftunnel up` | 启动隧道（自动下载 cloudflared） |
| `cftunnel down` | 停止隧道 |
| `cftunnel status` | 查看隧道状态 |
| `cftunnel logs [-f]` | 查看日志（`-f` 实时跟踪） |

### 系统服务

| 命令 | 说明 |
|------|------|
| `cftunnel install` | 注册系统服务（macOS launchd / Linux systemd） |
| `cftunnel uninstall` | 卸载系统服务 |

### 隧道生命周期

| 命令 | 说明 |
|------|------|
| `cftunnel destroy [--force]` | 删除隧道 + 所有 DNS 记录 |
| `cftunnel reset [--force]` | 完全重置（删隧道 + 清本地配置） |

### 版本管理

| 命令 | 说明 |
|------|------|
| `cftunnel version [--check]` | 显示版本 / 检查更新 |
| `cftunnel update` | 自动更新到最新版 |

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="scenarios">典型使用场景</h2>

### 场景 1: 暴露本地开发服务

```bash
cftunnel init
cftunnel create dev-tunnel
cftunnel add dev 3000 --domain dev.example.com
cftunnel up
# 现在 dev.example.com 指向 localhost:3000
```

### 场景 2: Webhook 接收

```bash
cftunnel add webhook 9801 --domain webhook.example.com
cftunnel up
cftunnel install  # 开机自启
```

### 场景 3: 多服务同时暴露

```bash
cftunnel add api 8080 --domain api.example.com
cftunnel add web 3000 --domain web.example.com
cftunnel add admin 8888 --domain admin.example.com
cftunnel list
```

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="config">配置文件</h2>

配置存储在 `~/.cftunnel/config.yml`：

```yaml
version: 1
auth:
  api_token: "your-token"
  account_id: "your-account-id"
tunnel:
  id: "tunnel-uuid"
  name: "my-tunnel"
  token: "tunnel-run-token"
routes:
  - name: myapp
    hostname: app.example.com
    service: http://localhost:3000
    zone_id: "auto-detected"
    dns_record_id: "auto-created"
```

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="ai">AI 助手集成</h2>

cftunnel 内置了 AI 助手 Skills，让 Claude Code、OpenClaw 等 AI 编码助手可以直接帮你管理隧道。

### Claude Code

将项目克隆到本地后，Claude Code 会自动加载 `.claude/skills/cftunnel.md`，你可以直接说：

```
帮我用 cftunnel quick 把本地 3000 端口临时分享出去
```

或者：

```
帮我用 cftunnel 把本地 3000 端口暴露到 dev.example.com
```

### OpenClaw / 其他 AI 助手

复制以下提示词给你的 AI 助手，即可让它帮你操作 cftunnel：

<details>
<summary>点击展开完整提示词</summary>

```
你是一个熟悉 cftunnel 的运维助手。cftunnel 是 Cloudflare Tunnel 的 CLI 管理工具，用于内网穿透。

## 安装

macOS / Linux:
curl -fsSL https://raw.githubusercontent.com/qingchencloud/cftunnel/main/install.sh | bash

Windows (PowerShell):
irm https://raw.githubusercontent.com/qingchencloud/cftunnel/main/install.ps1 | iex

## 两种模式

### 模式一：免域名（零配置，临时用）

cftunnel quick <端口>
# 自动生成 *.trycloudflare.com 随机域名，Ctrl+C 退出即失效
# 无需账户、Token、域名

### 模式二：自有域名（稳定持久）

需要引导用户准备两个参数：

1. API Token（API 令牌）:
   - 打开 https://dash.cloudflare.com/profile/api-tokens
   - 点击「创建令牌」→「创建自定义令牌」→「开始使用」
   - 添加 3 条权限（点「+ 添加更多」逐条添加）:
     帐户 | Cloudflare Tunnel | 编辑
     区域 | DNS              | 编辑（注意: 选「DNS」不是「DNS 设置」）
     区域 | 区域设置          | 读取
   - 第 2、3 行需先将左侧下拉框从「帐户」切换为「区域」
   - 区域资源 → 包括 → 特定区域 → 选择域名
   - 创建后立即复制令牌（只显示一次）

2. Account ID（账户 ID）:
   - 方式 A: https://dash.cloudflare.com → 点击域名 → 右下角「API」区域 → 账户 ID
   - 方式 B: 首页 → 账户名称旁「⋯」→ 复制账户 ID

使用流程：

# 第 1 步: 配置认证（需要上面两个参数）
cftunnel init --token <API_TOKEN> --account <ACCOUNT_ID>

# 第 2 步: 创建隧道
cftunnel create <隧道名称>

# 第 3 步: 添加路由（自动创建 DNS CNAME 记录）
cftunnel add <路由名称> <本地端口> --domain <完整域名>

# 第 4 步: 启动
cftunnel up

## 其他命令

- cftunnel quick <端口>   # 免域名模式（零配置）
- cftunnel down          # 停止隧道
- cftunnel status        # 查看状态
- cftunnel list          # 列出所有路由
- cftunnel remove <名称> # 删除路由（自动清理 DNS）
- cftunnel destroy       # 删除隧道 + 所有 DNS 记录
- cftunnel install       # 注册系统服务（开机自启）
- cftunnel logs -f       # 实时查看日志
- cftunnel update        # 自动更新 cftunnel

## 注意事项

- 临时分享优先推荐 `cftunnel quick`，零配置最快
- 自有域名模式需先完成 `init` 和 `create`
- 添加路由时会自动创建 DNS CNAME 记录，删除时自动清理
- 一个隧道可以挂载多条路由（多个域名指向不同本地端口）
- 域名必须是用户 Cloudflare 账户中已有的域名的子域名
```

</details>

<p align="right"><a href="#cftunnel">⬆ 回到顶部</a></p>

<h2 id="dev">开发</h2>

```bash
make build              # 本地构建
git tag v0.x.0 && git push --tags  # 推送 tag 自动触发 GitHub Actions 发版
```

<h2 id="contact">交流</h2>

- 官网: [cftunnel.qt.cool](https://cftunnel.qt.cool)
- QQ 群: [OpenClaw 交流群](https://qm.qq.com/q/qUfdR0jJVS)
- Issues: [GitHub Issues](https://github.com/qingchencloud/cftunnel/issues)

<h2 id="license">License</h2>

MIT

---

由 [武汉晴辰天下网络科技有限公司](https://qingchencloud.com) 开源维护
