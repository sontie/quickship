# QuickShip

基于 Go 的轻量级部署工具，通过 SSH Agent Forwarding 实现无密钥分发的并行化部署。

## 特性

- **SSH Agent Forwarding** - 转发本地私钥到远程服务器，无需在服务器存储密钥
- **并行部署** - 多主机并发执行，提升部署效率
- **实时日志** - 彩色输出远程执行日志，支持多主机区分
- **幂等操作** - Git 自动判断 clone 或 pull，可重复执行

## 安装

### 方式一：使用安装脚本（推荐）

**Linux/macOS:**
```bash
# 下载对应平台的二进制文件和安装脚本
chmod +x install.sh
./install.sh
```

**Windows:**
```cmd
# 下载 qship.exe 和 install.bat
install.bat
```

### 方式二：手动编译

需要 Go 1.19+ 环境：

```bash
go build -o qship
```

### 方式三：交叉编译多平台版本

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o qship-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o qship-darwin

# Windows
GOOS=windows GOARCH=amd64 go build -o qship.exe
```

## 快速开始

### 1. 初始化配置

```bash
qship init
```

生成 `deploy.yaml` 配置文件模板。

### 2. 编辑配置

编辑 `deploy.yaml`，配置主机和项目信息：

```yaml
version: "1.0"

hosts:
  ali-dev:
    addr: "192.168.1.100:22"
    user: "deploy"
  ali-prod:
    addr: "192.168.1.101:22"
    user: "deploy"

projects:
  - name: "api-server"
    repo: "git@github.com:user/api-server.git"
    path: "/opt/api-server"
    envs:
      dev: ["ali-dev"]
      prod: ["ali-prod"]
    scripts:
      deploy: |
        docker-compose down
        docker-compose up -d --build
```

### 3. 检查 SSH Agent

```bash
qship check
```

确保 ssh-agent 正在运行且已加载密钥。如果未启动：

```bash
eval $(ssh-agent)
ssh-add ~/.ssh/id_rsa
```

### 4. 查看配置

```bash
qship list
```

显示所有主机和项目的配置信息。

### 5. 执行部署

```bash
qship deploy dev
```

部署到指定环境（dev/prod）。

## 命令说明

### qship init
生成默认配置文件 `deploy.yaml`。

### qship check
检查 SSH Agent 状态和已加载的密钥。

### qship auth <host>
将本地公钥复制到远程主机（等同于 ssh-copy-id）。

```bash
qship auth deploy@192.168.1.100
```

### qship list
列出所有配置的主机和项目信息。

### qship deploy <env>
部署到指定环境。

```bash
qship deploy dev    # 部署到开发环境
qship deploy prod   # 部署到生产环境
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
