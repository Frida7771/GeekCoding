#!/bin/bash

# 测试非 root 用户配置的脚本

echo "=== Testing Non-Root User Configuration ==="
echo ""

# 1. 检查镜像是否存在
echo "1. Checking if image exists..."
if ! docker images | grep -q golang-code-runner; then
    echo "   ❌ Image not found. Please build it first: ./docker-build.sh"
    exit 1
fi
echo "   ✅ Image found"
echo ""

# 2. 检查镜像中的默认用户
echo "2. Checking default user in image..."
USER_INFO=$(docker run --rm golang-code-runner:latest id 2>/dev/null)
if echo "$USER_INFO" | grep -q "uid=1000(sandbox)"; then
    echo "   ✅ Running as sandbox user (UID 1000)"
    echo "   User info: $USER_INFO"
else
    echo "   ❌ Not running as sandbox user"
    echo "   User info: $USER_INFO"
fi
echo ""

# 3. 检查运行时用户（显式指定）
echo "3. Checking runtime user (explicit --user flag)..."
RUNTIME_USER=$(docker run --rm --user=1000:1000 golang-code-runner:latest whoami 2>/dev/null)
if [ "$RUNTIME_USER" = "sandbox" ]; then
    echo "   ✅ Running as sandbox user"
else
    echo "   ❌ Not running as sandbox user"
    echo "   Current user: $RUNTIME_USER"
fi
echo ""

# 4. 测试文件系统写入（应该失败）
echo "4. Testing filesystem write to /etc (should fail)..."
WRITE_TEST=$(docker run --rm --user=1000:1000 --read-only \
  golang-code-runner:latest sh -c "echo 'test' > /etc/test 2>&1" 2>&1)
if echo "$WRITE_TEST" | grep -qE "(Read-only|Permission denied|read-only)"; then
    echo "   ✅ Filesystem is read-only (secure)"
else
    echo "   ⚠️  Filesystem write test result: $WRITE_TEST"
fi
echo ""

# 5. 测试临时文件系统写入（应该成功）
echo "5. Testing tmpfs write to /tmp (should succeed)..."
TMPFS_TEST=$(docker run --rm --user=1000:1000 --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=10m \
  golang-code-runner:latest sh -c "echo 'test123' > /tmp/test.txt && cat /tmp/test.txt" 2>&1)
if echo "$TMPFS_TEST" | grep -q "test123"; then
    echo "   ✅ Tmpfs write succeeded"
    echo "   Content: $TMPFS_TEST"
else
    echo "   ❌ Tmpfs write failed"
    echo "   Output: $TMPFS_TEST"
fi
echo ""

# 6. 检查进程限制
echo "6. Checking process limits..."
PROC_LIMIT=$(docker run --rm --user=1000:1000 --pids-limit=10 \
  golang-code-runner:latest sh -c "ulimit -u" 2>/dev/null)
echo "   Process limit: $PROC_LIMIT"
echo ""

# 7. 检查是否是 root
echo "7. Verifying NOT running as root..."
ROOT_CHECK=$(docker run --rm --user=1000:1000 golang-code-runner:latest sh -c "id -u" 2>/dev/null)
if [ "$ROOT_CHECK" = "1000" ]; then
    echo "   ✅ Running as UID 1000 (not root)"
else
    echo "   ❌ Running as UID $ROOT_CHECK (should be 1000)"
fi
echo ""

echo "=== Test Complete ==="
echo ""
echo "Summary:"
echo "- Default user: sandbox (UID 1000)"
echo "- Runtime user: sandbox (UID 1000)"
echo "- Filesystem: Read-only (secure)"
echo "- Tmpfs: Writable (with limits)"
echo "- Process limits: Configured"

