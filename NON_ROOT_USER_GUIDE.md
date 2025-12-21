# 非 Root 用户使用指南

## 当前配置

系统已经配置为在 Docker 容器中以**非 root 用户**运行代码，提高安全性。

### 配置详情

1. **Dockerfile 中创建的用户**：
   - 用户名：`sandbox`
   - UID/GID：`1000:1000`
   - 组名：`sandbox`

2. **Docker 运行时的用户指定**：
   - `--user=1000:1000` 强制使用非 root 用户

---

## 验证非 Root 用户配置

### 1. 检查 Dockerfile 配置

查看 `internal/code/Dockerfile`：

```dockerfile
# 创建非 root 用户
RUN addgroup -g 1000 sandbox && \
    adduser -D -u 1000 -G sandbox sandbox

# 切换到非 root 用户
USER sandbox
```

### 2. 验证镜像中的用户

构建镜像后，可以检查：

```bash
# 构建镜像
./docker-build.sh

# 检查镜像中的用户
docker run --rm golang-code-runner:latest id
```

**预期输出**：
```
uid=1000(sandbox) gid=1000(sandbox) groups=1000(sandbox)
```

### 3. 验证运行时用户

运行一个测试容器，检查当前用户：

```bash
# 运行容器并检查用户
docker run --rm --user=1000:1000 golang-code-runner:latest id
```

**预期输出**：
```
uid=1000(sandbox) gid=1000(sandbox) groups=1000(sandbox)
```

### 4. 验证无法以 Root 运行

尝试以 root 运行（应该被阻止或需要显式指定）：

```bash
# 即使不指定 --user，由于 Dockerfile 中的 USER sandbox，也会以 sandbox 运行
docker run --rm golang-code-runner:latest whoami
```

**预期输出**：
```
sandbox
```

---

## 测试非 Root 用户限制

### 1. 测试文件系统写入权限

创建一个测试，尝试写入受保护的文件：

```bash
# 运行容器，尝试写入 /etc/passwd（应该失败）
docker run --rm --user=1000:1000 --read-only \
  golang-code-runner:latest sh -c "echo 'test' > /etc/passwd"
```

**预期结果**：
- 权限被拒绝（Permission denied）
- 或者文件系统只读错误

### 2. 测试特权操作

```bash
# 尝试执行需要 root 权限的操作
docker run --rm --user=1000:1000 \
  golang-code-runner:latest sh -c "mount -t tmpfs tmpfs /mnt"
```

**预期结果**：
- 操作失败（需要 root 权限）

### 3. 测试当前用户

在代码提交时，可以在用户代码中添加：

```go
package main
import (
    "fmt"
    "os"
    "os/user"
)

func main() {
    u, _ := user.Current()
    fmt.Println("Current user:", u.Username)
    fmt.Println("UID:", u.Uid)
    fmt.Println("GID:", u.Gid)
    
    // 尝试写入文件（应该失败或受限）
    err := os.WriteFile("/tmp/test.txt", []byte("test"), 0644)
    if err != nil {
        fmt.Println("Write failed:", err)
    }
}
```

**预期输出**：
```
Current user: sandbox
UID: 1000
GID: 1000
Write failed: ... (如果 /tmp 不可写)
```

---

## 在代码中检查用户

### Go 代码示例

```go
package main

import (
    "fmt"
    "os"
    "os/user"
)

func main() {
    // 获取当前用户
    u, err := user.Current()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Username: %s\n", u.Username)
    fmt.Printf("UID: %s\n", u.Uid)
    fmt.Printf("GID: %s\n", u.Gid)
    fmt.Printf("Home: %s\n", u.HomeDir)
    
    // 检查是否是 root
    if u.Uid == "0" {
        fmt.Println("WARNING: Running as root!")
    } else {
        fmt.Println("Running as non-root user (secure)")
    }
}
```

---

## 手动测试容器

### 1. 交互式进入容器

```bash
# 以非 root 用户进入容器
docker run --rm -it --user=1000:1000 \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=10m \
  golang-code-runner:latest sh
```

在容器内：

```bash
# 检查当前用户
whoami
# 输出: sandbox

# 检查用户 ID
id
# 输出: uid=1000(sandbox) gid=1000(sandbox) groups=1000(sandbox)

# 尝试写入根文件系统（应该失败）
echo "test" > /etc/test
# 错误: /etc/test: Read-only file system

# 尝试写入 /tmp（应该成功，但有大小限制）
echo "test" > /tmp/test.txt
cat /tmp/test.txt
```

### 2. 测试二进制文件执行

```bash
# 挂载一个测试二进制文件
docker run --rm --user=1000:1000 \
  -v /path/to/binary:/app/runner:ro \
  golang-code-runner:latest /app/runner
```

---

## 验证安全配置

### 完整的安全测试脚本

创建 `test-sandbox.sh`：

```bash
#!/bin/bash

echo "=== Testing Non-Root User Configuration ==="

# 1. 检查镜像中的用户
echo "1. Checking user in image..."
docker run --rm golang-code-runner:latest id

# 2. 检查运行时用户
echo "2. Checking runtime user..."
docker run --rm --user=1000:1000 golang-code-runner:latest whoami

# 3. 测试文件系统写入（应该失败）
echo "3. Testing filesystem write (should fail)..."
docker run --rm --user=1000:1000 --read-only \
  golang-code-runner:latest sh -c "echo 'test' > /etc/test" 2>&1 | head -1

# 4. 测试临时文件系统写入（应该成功）
echo "4. Testing tmpfs write (should succeed)..."
docker run --rm --user=1000:1000 --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=10m \
  golang-code-runner:latest sh -c "echo 'test' > /tmp/test.txt && cat /tmp/test.txt"

# 5. 检查进程限制
echo "5. Checking process limits..."
docker run --rm --user=1000:1000 --pids-limit=10 \
  golang-code-runner:latest sh -c "ulimit -u"

echo "=== Test Complete ==="
```

运行测试：

```bash
chmod +x test-sandbox.sh
./test-sandbox.sh
```

---

## 常见问题

### Q1: 如何确认容器以非 root 运行？

**A**: 在代码中添加用户检查，或运行：
```bash
docker run --rm golang-code-runner:latest id
```

### Q2: 如果用户代码需要写入文件怎么办？

**A**: 使用 `/tmp` 目录（通过 `--tmpfs` 挂载），但注意大小限制（10MB）。

### Q3: 如何修改用户 ID？

**A**: 修改 `Dockerfile` 中的 UID/GID 和 `service/submit.go` 中的 `--user` 参数。

### Q4: 为什么需要非 root 用户？

**A**: 
- 防止权限提升攻击
- 限制文件系统访问
- 提高容器安全性
- 符合安全最佳实践

### Q5: 如何调试权限问题？

**A**: 
1. 检查容器内的用户：`docker run --rm golang-code-runner:latest id`
2. 检查文件权限：`docker run --rm golang-code-runner:latest ls -la /app`
3. 查看错误日志：检查 `stderr` 输出

---

## 生产环境验证清单

- [ ] 镜像构建成功
- [ ] 容器以非 root 用户运行（UID 1000）
- [ ] 文件系统为只读（除了 /tmp）
- [ ] 无法执行特权操作
- [ ] 进程数限制生效
- [ ] CPU 限制生效
- [ ] 内存限制生效
- [ ] 网络隔离生效

---

## 总结

系统已经配置为：
- ✅ 在 Dockerfile 中创建非 root 用户 `sandbox` (UID 1000)
- ✅ 运行时强制使用 `--user=1000:1000`
- ✅ 文件系统只读（除了临时目录）
- ✅ 所有 capabilities 被删除
- ✅ 禁止获取新权限

这确保了代码在安全的沙盒环境中运行，无法进行权限提升或系统破坏。

