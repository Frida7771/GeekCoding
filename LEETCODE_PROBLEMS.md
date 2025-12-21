# LeetCode 标准题目数据

## 使用说明

这些题目可以直接通过 API 或 SQL 插入到数据库中。

### API 方式（推荐）

使用 `/admin/problem-create` 接口，需要管理员权限。

### SQL 方式

可以直接执行 SQL 语句插入数据。

---

## 题目列表

### 1. 两数之和 (Two Sum)

**题目信息**:
- **Title**: 两数之和
- **Content**: 
```
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那 两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。

你可以按任意顺序返回答案。

示例 1：
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

示例 2：
输入：nums = [3,2,4], target = 6
输出：[1,2]

示例 3：
输入：nums = [3,3], target = 6
输出：[0,1]

提示：
- 2 <= nums.length <= 10^4
- -10^9 <= nums[i] <= 10^9
- -10^9 <= target <= 10^9
- 只会存在一个有效答案
```
- **Max Runtime**: 3000 (3秒)
- **Max Memory**: 64 (MB)
- **Category**: 数组、哈希表

**测试用例**:
```json
[
  {"input": "2 7 11 15\n9", "output": "[0,1]"},
  {"input": "3 2 4\n6", "output": "[1,2]"},
  {"input": "3 3\n6", "output": "[0,1]"},
  {"input": "2 5 5 11\n10", "output": "[1,2]"}
]
```

**Go 代码示例**:
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
    
    // 实现两数之和逻辑
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            if nums[i] + nums[j] == target {
                fmt.Printf("[%d,%d]", i, j)
                return
            }
        }
    }
}
```

---

### 2. 反转链表 (Reverse Linked List)

**题目信息**:
- **Title**: 反转链表
- **Content**:
```
给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。

示例 1：
输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]

示例 2：
输入：head = [1,2]
输出：[2,1]

示例 3：
输入：head = []
输出：[]

提示：
- 链表中节点的数目范围是 [0, 5000]
- -5000 <= Node.val <= 5000
```
- **Max Runtime**: 3000
- **Max Memory**: 64
- **Category**: 链表、递归

**测试用例**:
```json
[
  {"input": "5\n1 2 3 4 5", "output": "5 4 3 2 1"},
  {"input": "2\n1 2", "output": "2 1"},
  {"input": "0", "output": ""},
  {"input": "1\n1", "output": "1"}
]
```

---

### 3. 有效的括号 (Valid Parentheses)

**题目信息**:
- **Title**: 有效的括号
- **Content**:
```
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：
1. 左括号必须用相同类型的右括号闭合。
2. 左括号必须以正确的顺序闭合。
3. 每个右括号都有一个对应的相同类型的左括号。

示例 1：
输入：s = "()"
输出：true

示例 2：
输入：s = "()[]{}"
输出：true

示例 3：
输入：s = "(]"
输出：false

示例 4：
输入：s = "([)]"
输出：false

示例 5：
输入：s = "{[]}"
输出：true
```
- **Max Runtime**: 2000
- **Max Memory**: 32
- **Category**: 栈、字符串

**测试用例**:
```json
[
  {"input": "()", "output": "true"},
  {"input": "()[]{}", "output": "true"},
  {"input": "(]", "output": "false"},
  {"input": "([)]", "output": "false"},
  {"input": "{[]}", "output": "true"},
  {"input": "(((", "output": "false"}
]
```

---

### 4. 合并两个有序链表 (Merge Two Sorted Lists)

**题目信息**:
- **Title**: 合并两个有序链表
- **Content**:
```
将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。

示例 1：
输入：l1 = [1,2,4], l2 = [1,3,4]
输出：[1,1,2,3,4,4]

示例 2：
输入：l1 = [], l2 = []
输出：[]

示例 3：
输入：l1 = [], l2 = [0]
输出：[0]

提示：
- 两个链表的节点数目范围是 [0, 50]
- -100 <= Node.val <= 100
- l1 和 l2 均按 非递减顺序 排列
```
- **Max Runtime**: 3000
- **Max Memory**: 64
- **Category**: 链表、递归

**测试用例**:
```json
[
  {"input": "3\n1 2 4\n3\n1 3 4", "output": "1 1 2 3 4 4"},
  {"input": "0\n0", "output": ""},
  {"input": "0\n1\n0", "output": "0"},
  {"input": "1\n1\n1\n2", "output": "1 2"}
]
```

---

### 5. 最大子数组和 (Maximum Subarray)

**题目信息**:
- **Title**: 最大子数组和
- **Content**:
```
给你一个整数数组 nums ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

子数组 是数组中的一个连续部分。

示例 1：
输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
输出：6
解释：连续子数组 [4,-1,2,1] 的和最大，为 6 。

示例 2：
输入：nums = [1]
输出：1

示例 3：
输入：nums = [5,4,-1,7,8]
输出：23

提示：
- 1 <= nums.length <= 10^5
- -10^4 <= nums[i] <= 10^4
```
- **Max Runtime**: 5000
- **Max Memory**: 64
- **Category**: 数组、动态规划

**测试用例**:
```json
[
  {"input": "9\n-2 1 -3 4 -1 2 1 -5 4", "output": "6"},
  {"input": "1\n1", "output": "1"},
  {"input": "5\n5 4 -1 7 8", "output": "23"},
  {"input": "3\n-1 -2 -3", "output": "-1"}
]
```

---

### 6. 爬楼梯 (Climbing Stairs)

**题目信息**:
- **Title**: 爬楼梯
- **Content**:
```
假设你正在爬楼梯。需要 n 阶你才能到达楼顶。

每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？

示例 1：
输入：n = 2
输出：2
解释：有两种方法可以爬到楼顶。
1. 1 阶 + 1 阶
2. 2 阶

示例 2：
输入：n = 3
输出：3
解释：有三种方法可以爬到楼顶。
1. 1 阶 + 1 阶 + 1 阶
2. 1 阶 + 2 阶
3. 2 阶 + 1 阶

提示：
- 1 <= n <= 45
```
- **Max Runtime**: 2000
- **Max Memory**: 32
- **Category**: 动态规划、数学

**测试用例**:
```json
[
  {"input": "2", "output": "2"},
  {"input": "3", "output": "3"},
  {"input": "4", "output": "5"},
  {"input": "5", "output": "8"},
  {"input": "1", "output": "1"}
]
```

---

### 7. 买卖股票的最佳时机 (Best Time to Buy and Sell Stock)

**题目信息**:
- **Title**: 买卖股票的最佳时机
- **Content**:
```
给定一个数组 prices ，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。

你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。设计一个算法来计算你所能获取的最大利润。

返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 0 。

示例 1：
输入：[7,1,5,3,6,4]
输出：5
解释：在第 2 天（股票价格 = 1）的时候买入，在第 5 天（股票价格 = 6）的时候卖出，最大利润 = 6-1 = 5 。
     注意利润不能是 7-1 = 6, 因为卖出价格需要大于买入价格；同时，你不能在买入前卖出股票。

示例 2：
输入：prices = [7,6,4,3,1]
输出：0
解释：在这种情况下, 交易无法完成, 所以最大利润为 0。

提示：
- 1 <= prices.length <= 10^5
- 0 <= prices[i] <= 10^4
```
- **Max Runtime**: 5000
- **Max Memory**: 64
- **Category**: 数组、动态规划

**测试用例**:
```json
[
  {"input": "6\n7 1 5 3 6 4", "output": "5"},
  {"input": "5\n7 6 4 3 1", "output": "0"},
  {"input": "2\n1 2", "output": "1"},
  {"input": "3\n2 4 1", "output": "2"}
]
```

---

### 8. 对称二叉树 (Symmetric Tree)

**题目信息**:
- **Title**: 对称二叉树
- **Content**:
```
给你一个二叉树的根节点 root ， 检查它是否轴对称。

示例 1：
输入：root = [1,2,2,3,4,4,3]
输出：true

示例 2：
输入：root = [1,2,2,null,3,null,3]
输出：false

提示：
- 树中节点数目在范围 [1, 1000] 内
- -100 <= Node.val <= 100
```
- **Max Runtime**: 3000
- **Max Memory**: 64
- **Category**: 树、递归

**测试用例**:
```json
[
  {"input": "7\n1 2 2 3 4 4 3", "output": "true"},
  {"input": "7\n1 2 2 -1 3 -1 3", "output": "false"},
  {"input": "1\n1", "output": "true"},
  {"input": "3\n1 2 2", "output": "true"}
]
```

---

### 9. 二叉树的最大深度 (Maximum Depth of Binary Tree)

**题目信息**:
- **Title**: 二叉树的最大深度
- **Content**:
```
给定一个二叉树，找出其最大深度。

二叉树的深度为根节点到最远叶子节点的最长路径上的节点数。

说明: 叶子节点是指没有子节点的节点。

示例：
给定二叉树 [3,9,20,null,null,15,7]，

    3
   / \
  9  20
    /  \
   15   7

返回它的最大深度 3 。
```
- **Max Runtime**: 3000
- **Max Memory**: 64
- **Category**: 树、深度优先搜索

**测试用例**:
```json
[
  {"input": "7\n3 9 20 -1 -1 15 7", "output": "3"},
  {"input": "1\n1", "output": "1"},
  {"input": "3\n1 2 3", "output": "2"},
  {"input": "5\n1 -1 2 -1 3", "output": "3"}
]
```

---

### 10. 回文数 (Palindrome Number)

**题目信息**:
- **Title**: 回文数
- **Content**:
```
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。

回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。

例如，121 是回文，而 123 不是。

示例 1：
输入：x = 121
输出：true

示例 2：
输入：x = -121
输出：false
解释：从左向右读, 为 -121 。 从右向左读, 为 121- 。因此它不是一个回文数。

示例 3：
输入：x = 10
输出：false
解释：从右向左读, 为 01 。因此它不是一个回文数。

提示：
- -2^31 <= x <= 2^31 - 1
```
- **Max Runtime**: 2000
- **Max Memory**: 32
- **Category**: 数学

**测试用例**:
```json
[
  {"input": "121", "output": "true"},
  {"input": "-121", "output": "false"},
  {"input": "10", "output": "false"},
  {"input": "0", "output": "true"},
  {"input": "1221", "output": "true"},
  {"input": "123", "output": "false"}
]
```

---

## API 调用示例

### 使用 curl 创建题目

```bash
# 1. 两数之和
curl -X POST http://localhost:8080/admin/problem-create \
  -H "Authorization: <admin-token>" \
  -F "title=两数之和" \
  -F "content=给定一个整数数组 nums 和一个整数目标值 target..." \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F "test_cases={\"input\":\"2 7 11 15\n9\",\"output\":\"[0,1]\"}" \
  -F "test_cases={\"input\":\"3 2 4\n6\",\"output\":\"[1,2]\"}" \
  -F "test_cases={\"input\":\"3 3\n6\",\"output\":\"[0,1]\"}"
```

---

## SQL 插入语句（备用）

如果需要直接使用 SQL 插入，需要先创建分类，然后插入题目和测试用例。

**注意**: 这些 SQL 语句需要根据实际的分类 ID 进行调整。

