# QuickShip

基于 Go 的轻量级部署工具，通过 SSH Agent Forwarding 实现无密钥分发的并行化部署。

## 特性

- **SSH Agent Forwarding** - 转发本地私钥到远程服务器，无需在服务器存储密钥
- **并行部署** - 多主机并发执行，提升部署效率
- **实时日志** - 彩色输出远程执行日志，支持多主机区分
- **幂等操作** - Git 自动判断 clone 或 pull，可重复执行

## 安装

### 方式一：一键安装（推荐）

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/install.ps1 | iex
```

### 卸载

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/uninstall.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/uninstall.ps1 | iex
```

或手动删除：
- Linux/macOS: `sudo rm /usr/local/bin/qship && rm -rf ~/.quickship`
- Windows: 删除 `qship.exe` 和 `%USERPROFILE%\.quickship` 目录

### 升级

**如果已安装 v0.1.1+ 或更新版本：**
```bash
qship upgrade
```

**如果安装的是旧版本（无 upgrade 命令）：**

重新运行安装脚本会自动覆盖旧版本：

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/install.ps1 | iex
```

或手动下载最新二进制文件替换：
```bash
# 从 GitHub Releases 获取最新版本
curl -L -o qship https://github.com/sontie/quickship/releases/latest/download/qship-darwin-arm64
chmod +x qship
sudo mv qship /usr/local/bin/qship
```

### 方式二：下载二进制文件

从 [Releases](https://github.com/sontie/quickship/releases) 下载对应平台的文件：

**Linux/macOS:**
```bash
chmod +x qship-*
sudo mv qship-* /usr/local/bin/qship
```

**Windows:**
下载 `qship-windows-amd64.exe`，重命名为 `qship.exe` 并添加到 PATH。

### 方式三：从源码编译

需要 Go 1.19+ 环境：

```bash
git clone https://github.com/sontie/quickship.git
cd quickship
go build -o qship
```

## 快速开始

### 1. 初始化配置

```bash
qship init
```

生成 `qship.yaml` 配置文件模板。

### 2. 编辑配置

编辑 `qship.yaml`，配置主机和项目信息。

#### 配置文件说明

```yaml
# 配置文件版本
version: "1.0"

# 主机列表
hosts:
  ali-dev:                    # 主机别名
    addr: "192.168.1.100:22"  # SSH 地址和端口
    user: "deploy"            # SSH 用户名
  ali-prod:
    addr: "192.168.1.101:22"
    user: "deploy"

# 项目列表
projects:
  - name: "api-server"                        # 项目名称
    repo: "git@github.com:user/api-server.git" # Git 仓库地址
    path: "/opt/api-server"                   # 远程服务器部署路径
    envs:                                     # 环境配置
      dev: ["ali-dev"]                        # dev 环境部署到 ali-dev 主机
      prod: ["ali-prod"]                      # prod 环境部署到 ali-prod 主机
    scripts:                                  # 部署脚本
      deploy: |                               # deploy 脚本内容
        docker-compose down
        docker-compose up -d --build
```

### 3. 环境准备

运行以下命令检查 SSH 环境，如有问题会给出引导：

```bash
qship check
```

如果你已熟悉 SSH，可自行配置；如果不熟悉，按照提示操作即可。

### 4. 查看配置

```bash
qship list
# 或使用简写
qship ls
```

显示所有主机和项目的配置信息。

### 5. 执行部署

```bash
qship go dev
```

部署到指定环境（dev/prod）。

## 命令说明

### qship version / -v / --version
查看版本信息。

```bash
qship version
qship -v
qship --version
```

### qship init
生成默认配置文件 `qship.yaml`。

### qship check
检查 SSH Agent 状态和已加载的密钥。

### qship auth <host>
将本地公钥复制到远程主机（等同于 ssh-copy-id）。

```bash
qship auth deploy@192.168.1.100
```

### qship list / ls
列出所有配置的主机和项目信息。

### qship go <env>
部署到指定环境。

```bash
qship go dev    # 部署到开发环境
qship go prod   # 部署到生产环境
```

### qship upgrade
自动检查并升级到最新版本。

```bash
qship upgrade
```

如果安装在系统目录需要权限：
```bash
sudo qship upgrade
```

### qship exec "<command>" [hosts]
在指定主机上执行命令。

```bash
qship exec "uptime"                    # 在所有主机执行
qship exec "df -h" ali-dev,ali-prod   # 在指定主机执行
```

## 工作原理

1. 从本地 `SSH_AUTH_SOCK` 获取 ssh-agent 连接
2. 使用 Agent Forwarding 建立 SSH 连接
3. 在远程服务器执行 Git 操作（clone/pull）
4. 执行配置的部署脚本
5. 实时回显远程执行日志

## License

MIT
