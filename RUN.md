# GeekCoding 运行指南

## 前置要求

1. **Go 1.21+**
2. **MySQL 8.0+**（运行在 `localhost:3306`）
3. **Redis**（运行在 `localhost:6379`）
4. **Docker**（用于运行用户提交的代码）

## 运行步骤

### 1. 启动 MySQL 和 Redis

确保 MySQL 和 Redis 服务正在运行：

```bash
# 检查 MySQL
mysql -u root -p'#Etnlhy1396917302' -e "SELECT 1;" Geek_Coding

# 检查 Redis
redis-cli ping
```

如果服务未运行，请启动它们。

### 2. 构建 Docker 镜像

构建用于运行用户代码的 Docker 镜像：

```bash
./docker-build.sh
```

或者手动构建：

```bash
docker build -t golang-code-runner:latest -f internal/code/Dockerfile .
```

验证镜像是否构建成功：

```bash
docker images | grep golang-code-runner
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 生成 Swagger 文档（可选）

如果需要更新 API 文档：

```bash
swag init
```

### 5. 运行应用

```bash
go run main.go
```

应用将在 `http://localhost:8080` 启动。

### 6. 访问 Swagger 文档

打开浏览器访问：

```
http://localhost:8080/swagger/index.html
```

## 验证运行

### 检查服务状态

1. **应用日志**：查看控制台输出，应该看到：
   ```
   [GIN-debug] Listening and serving HTTP on :8080
   ```

2. **Swagger 文档**：访问 `http://localhost:8080/swagger/index.html`，应该能看到所有 API 接口

3. **健康检查**：可以尝试访问：
   ```
   http://localhost:8080/problem-list
   ```

## 常见问题

### 1. MySQL 连接失败

**错误**：`gorm Init Error`

**解决**：
- 检查 MySQL 是否运行：`mysql -u root -p`
- 检查数据库是否存在：`SHOW DATABASES;`
- 检查连接信息是否正确（`models/init.go`）

### 2. Redis 连接失败

**错误**：Redis 连接超时

**解决**：
- 检查 Redis 是否运行：`redis-cli ping`
- 如果未运行，启动 Redis：`redis-server`

### 3. Docker 镜像未找到

**错误**：`docker: Error response from daemon: pull access denied for golang-code-runner`

**解决**：
- 确保已构建镜像：`./docker-build.sh`
- 检查镜像是否存在：`docker images | grep golang-code-runner`

### 4. 端口被占用

**错误**：`bind: address already in use`

**解决**：
- 查找占用端口的进程：`lsof -i :8080`
- 杀死进程或修改 `main.go` 中的端口

### 5. Swagger 文档为空

**解决**：
- 重新生成文档：`swag init`
- 确保所有 API 函数都有正确的 Swagger 注释

## 开发模式

### 热重载（使用 air）

安装 air：

```bash
go install github.com/cosmtrek/air@latest
```

创建 `.air.toml` 配置文件，然后运行：

```bash
air
```

## 生产部署

1. 编译二进制文件：
   ```bash
   go build -o GeekCoding main.go
   ```

2. 运行：
   ```bash
   ./GeekCoding
   ```

3. 或使用 systemd 等服务管理器管理进程

## 环境变量（可选）

可以创建 `.env` 文件来配置：

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=#Etnlhy1396917302
DB_NAME=Geek_Coding
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

然后在代码中使用 `os.Getenv()` 读取。

