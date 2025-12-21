# LeetCode 题目一键插入指南

## 快速开始

### 1. 获取管理员 Token

```bash
# 登录获取 token
curl -X POST http://localhost:8080/login \
  -d "username=your_admin_username" \
  -d "password=your_admin_password"
```

从响应中获取 `token` 字段的值。

### 2. 运行插入脚本

```bash
./insert_leetcode_problems.sh <your-admin-token>
```

例如：
```bash
./insert_leetcode_problems.sh "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 包含的10个题目

1. ✅ **两数之和** (Two Sum) - 数组、哈希表
2. ✅ **反转链表** (Reverse Linked List) - 链表、递归
3. ✅ **有效的括号** (Valid Parentheses) - 栈、字符串
4. ✅ **合并两个有序链表** (Merge Two Sorted Lists) - 链表
5. ✅ **最大子数组和** (Maximum Subarray) - 数组、动态规划
6. ✅ **爬楼梯** (Climbing Stairs) - 动态规划、数学
7. ✅ **买卖股票的最佳时机** (Best Time to Buy and Sell Stock) - 数组、动态规划
8. ✅ **对称二叉树** (Symmetric Tree) - 树、递归
9. ✅ **二叉树的最大深度** (Maximum Depth of Binary Tree) - 树、DFS
10. ✅ **回文数** (Palindrome Number) - 数学

## 注意事项

### 1. 分类 ID

脚本中使用的 `category_ids=1` 是默认值。如果您的系统中分类 ID 不同，需要：

1. 先创建分类：
```bash
curl -X POST http://localhost:8080/admin/category-create \
  -H "Authorization: <token>" \
  -F "name=算法"
```

2. 查看分类列表获取 ID：
```bash
curl -X GET "http://localhost:8080/admin/category-list" \
  -H "Authorization: <token>"
```

3. 修改脚本中的 `category_ids` 参数。

### 2. 测试用例格式

测试用例的输入格式需要适配 Go 代码的读取方式。当前格式：

**两数之和示例**:
```
输入: "4\n2 7 11 15\n9"
解释: 
  - 第一行: 数组长度 4
  - 第二行: 数组元素 2 7 11 15
  - 第三行: 目标值 9
输出: "[0,1]"
```

**Go 代码读取示例**:
```go
package main
import "fmt"
func main() {
    var n, target int
    fmt.Scan(&n)
    nums := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Scan(&nums[i])
    }
    fmt.Scan(&target)
    // ... 处理逻辑
    fmt.Println("[0,1]")
}
```

### 3. 输出格式

当前使用精确字符串匹配。输出必须与期望输出完全一致（包括格式）。

**建议**: 如果遇到格式问题，可以考虑：
- 去除首尾空白字符
- 统一换行符
- 参考 `TEST_CASE_LOGIC_ANALYSIS.md` 中的改进方案

## 手动插入单个题目

如果脚本执行失败，可以手动插入：

```bash
curl -X POST http://localhost:8080/admin/problem-create \
  -H "Authorization: <your-token>" \
  -F "title=两数之和" \
  -F "content=题目描述..." \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"4\n2 7 11 15\n9","output":"[0,1]"}' \
  -F 'test_cases={"input":"3\n3 2 4\n6","output":"[1,2]"}'
```

## 验证插入结果

```bash
# 查看题目列表
curl -X GET "http://localhost:8080/problem-list?page=1&size=10"

# 查看题目详情
curl -X GET "http://localhost:8080/problem-detail?identity=<problem-identity>"
```

## 题目详情

详细的题目描述、测试用例和代码示例请参考 `LEETCODE_PROBLEMS.md`。

## 故障排除

### 1. 401 Unauthorized
- 检查 token 是否有效
- 确认 token 是否过期（24小时有效期）
- 重新登录获取新 token

### 2. 分类不存在
- 先创建分类
- 或修改脚本使用已存在的分类 ID

### 3. 测试用例格式错误
- 检查 JSON 格式是否正确
- 确认输入输出格式是否匹配代码逻辑

### 4. 服务器未运行
- 确认服务器运行在 `http://localhost:8080`
- 检查服务器日志

