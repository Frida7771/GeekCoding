#!/bin/bash

# LeetCode 标准题目一键插入脚本
# 使用方法: ./insert_leetcode_problems.sh <admin-token>
# 需要先登录获取管理员 token

API_BASE="http://localhost:8080"
TOKEN=$1

if [ -z "$TOKEN" ]; then
    echo "错误: 请提供管理员 token"
    echo "使用方法: ./insert_leetcode_problems.sh <admin-token>"
    echo "获取 token: curl -X POST $API_BASE/login -d 'username=admin&password=password'"
    exit 1
fi

echo "开始插入 LeetCode 题目..."

# 1. 两数之和
echo "插入题目 1/10: 两数之和"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=两数之和" \
  -F "content=给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那 两个 整数，并返回它们的数组下标。

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
- 只会存在一个有效答案" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"4\n2 7 11 15\n9","output":"[0,1]"}' \
  -F 'test_cases={"input":"3\n3 2 4\n6","output":"[1,2]"}' \
  -F 'test_cases={"input":"2\n3 3\n6","output":"[0,1]"}' \
  -F 'test_cases={"input":"4\n2 5 5 11\n10","output":"[1,2]"}'

echo -e "\n"

# 2. 反转链表
echo "插入题目 2/10: 反转链表"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=反转链表" \
  -F "content=给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。

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
- -5000 <= Node.val <= 5000" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"5\n1 2 3 4 5","output":"5 4 3 2 1"}' \
  -F 'test_cases={"input":"2\n1 2","output":"2 1"}' \
  -F 'test_cases={"input":"0","output":""}' \
  -F 'test_cases={"input":"1\n1","output":"1"}'

echo -e "\n"

# 3. 有效的括号
echo "插入题目 3/10: 有效的括号"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=有效的括号" \
  -F "content=给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：
1. 左括号必须用相同类型的右括号闭合。
2. 左括号必须以正确的顺序闭合。
3. 每个右括号都有一个对应的相同类型的左括号。

示例 1：
输入：s = \"()\"
输出：true

示例 2：
输入：s = \"()[]{}\"
输出：true

示例 3：
输入：s = \"(]\"
输出：false

示例 4：
输入：s = \"([)]\"
输出：false

示例 5：
输入：s = \"{[]}\"
输出：true" \
  -F "max_runtime=2000" \
  -F "max_mem=32" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"()","output":"true"}' \
  -F 'test_cases={"input":"()[]{}","output":"true"}' \
  -F 'test_cases={"input":"(]","output":"false"}' \
  -F 'test_cases={"input":"([)]","output":"false"}' \
  -F 'test_cases={"input":"{[]}","output":"true"}' \
  -F 'test_cases={"input":"(((","output":"false"}'

echo -e "\n"

# 4. 合并两个有序链表
echo "插入题目 4/10: 合并两个有序链表"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=合并两个有序链表" \
  -F "content=将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。

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
- l1 和 l2 均按 非递减顺序 排列" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"3\n1 2 4\n3\n1 3 4","output":"1 1 2 3 4 4"}' \
  -F 'test_cases={"input":"0\n0","output":""}' \
  -F 'test_cases={"input":"0\n1\n0","output":"0"}' \
  -F 'test_cases={"input":"1\n1\n1\n2","output":"1 2"}'

echo -e "\n"

# 5. 最大子数组和
echo "插入题目 5/10: 最大子数组和"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=最大子数组和" \
  -F "content=给你一个整数数组 nums ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

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
- -10^4 <= nums[i] <= 10^4" \
  -F "max_runtime=5000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"9\n-2 1 -3 4 -1 2 1 -5 4","output":"6"}' \
  -F 'test_cases={"input":"1\n1","output":"1"}' \
  -F 'test_cases={"input":"5\n5 4 -1 7 8","output":"23"}' \
  -F 'test_cases={"input":"3\n-1 -2 -3","output":"-1"}'

echo -e "\n"

# 6. 爬楼梯
echo "插入题目 6/10: 爬楼梯"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=爬楼梯" \
  -F "content=假设你正在爬楼梯。需要 n 阶你才能到达楼顶。

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
- 1 <= n <= 45" \
  -F "max_runtime=2000" \
  -F "max_mem=32" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"2","output":"2"}' \
  -F 'test_cases={"input":"3","output":"3"}' \
  -F 'test_cases={"input":"4","output":"5"}' \
  -F 'test_cases={"input":"5","output":"8"}' \
  -F 'test_cases={"input":"1","output":"1"}'

echo -e "\n"

# 7. 买卖股票的最佳时机
echo "插入题目 7/10: 买卖股票的最佳时机"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=买卖股票的最佳时机" \
  -F "content=给定一个数组 prices ，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。

你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。设计一个算法来计算你所能获取的最大利润。

返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 0 。

示例 1：
输入：[7,1,5,3,6,4]
输出：5
解释：在第 2 天（股票价格 = 1）的时候买入，在第 5 天（股票价格 = 6）的时候卖出，最大利润 = 6-1 = 5 。

示例 2：
输入：prices = [7,6,4,3,1]
输出：0
解释：在这种情况下, 交易无法完成, 所以最大利润为 0。

提示：
- 1 <= prices.length <= 10^5
- 0 <= prices[i] <= 10^4" \
  -F "max_runtime=5000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"6\n7 1 5 3 6 4","output":"5"}' \
  -F 'test_cases={"input":"5\n7 6 4 3 1","output":"0"}' \
  -F 'test_cases={"input":"2\n1 2","output":"1"}' \
  -F 'test_cases={"input":"3\n2 4 1","output":"2"}'

echo -e "\n"

# 8. 对称二叉树
echo "插入题目 8/10: 对称二叉树"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=对称二叉树" \
  -F "content=给你一个二叉树的根节点 root ， 检查它是否轴对称。

示例 1：
输入：root = [1,2,2,3,4,4,3]
输出：true

示例 2：
输入：root = [1,2,2,null,3,null,3]
输出：false

提示：
- 树中节点数目在范围 [1, 1000] 内
- -100 <= Node.val <= 100" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"7\n1 2 2 3 4 4 3","output":"true"}' \
  -F 'test_cases={"input":"7\n1 2 2 -1 3 -1 3","output":"false"}' \
  -F 'test_cases={"input":"1\n1","output":"true"}' \
  -F 'test_cases={"input":"3\n1 2 2","output":"true"}'

echo -e "\n"

# 9. 二叉树的最大深度
echo "插入题目 9/10: 二叉树的最大深度"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=二叉树的最大深度" \
  -F "content=给定一个二叉树，找出其最大深度。

二叉树的深度为根节点到最远叶子节点的最长路径上的节点数。

说明: 叶子节点是指没有子节点的节点。

示例：
给定二叉树 [3,9,20,null,null,15,7]，

    3
   / \
  9  20
    /  \
   15   7

返回它的最大深度 3 。" \
  -F "max_runtime=3000" \
  -F "max_mem=64" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"7\n3 9 20 -1 -1 15 7","output":"3"}' \
  -F 'test_cases={"input":"1\n1","output":"1"}' \
  -F 'test_cases={"input":"3\n1 2 3","output":"2"}' \
  -F 'test_cases={"input":"5\n1 -1 2 -1 3","output":"3"}'

echo -e "\n"

# 10. 回文数
echo "插入题目 10/10: 回文数"
curl -X POST "$API_BASE/admin/problem-create" \
  -H "Authorization: $TOKEN" \
  -F "title=回文数" \
  -F "content=给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。

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
- -2^31 <= x <= 2^31 - 1" \
  -F "max_runtime=2000" \
  -F "max_mem=32" \
  -F "category_ids=1" \
  -F 'test_cases={"input":"121","output":"true"}' \
  -F 'test_cases={"input":"-121","output":"false"}' \
  -F 'test_cases={"input":"10","output":"false"}' \
  -F 'test_cases={"input":"0","output":"true"}' \
  -F 'test_cases={"input":"1221","output":"true"}' \
  -F 'test_cases={"input":"123","output":"false"}'

echo -e "\n"
echo "所有题目插入完成！"

