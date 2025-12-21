# Test Case 判断逻辑分析

## 当前实现

### 数据结构

```go
type TestCase struct {
    Input  string  // 输入数据
    Output string  // 期望输出（标准答案）
}
```

### 判断逻辑

**位置**: `service/submit.go` 第 345 行

```go
//答案错误
if tc.Output != out.String() {
    msg = "答案错误"
    WA <- 1
    return
}
```

## 工作流程

1. **输入阶段**:
   - 将 `testCase.Input` 通过 stdin 输入到用户代码
   ```go
   io.WriteString(stdinPipe, tc.Input)
   ```

2. **执行阶段**:
   - 用户代码在 Docker 容器中执行
   - 读取 stdin 输入
   - 处理并输出结果到 stdout

3. **输出捕获**:
   - 捕获用户代码的 stdout 输出
   ```go
   var out bytes.Buffer
   dockerCmd.Stdout = &out
   ```

4. **比较阶段**:
   - 将用户代码的实际输出 `out.String()` 与期望输出 `tc.Output` 进行**精确字符串比较**
   ```go
   if tc.Output != out.String() {
       // 答案错误
   }
   ```

## 当前逻辑的特点

### ✅ 优点
- **简单直接**: 字符串精确比较，逻辑清晰
- **准确**: 完全匹配才通过，不会误判

### ⚠️ 潜在问题

#### 1. **空白字符敏感**
- **问题**: 如果用户输出有尾随空格或换行符，但标准答案没有（或反之），会判错
- **示例**:
  ```
  标准答案: "42"
  用户输出: "42\n"  // 多了一个换行符
  结果: WA (错误)
  ```

#### 2. **换行符不一致**
- **问题**: Windows (`\r\n`) vs Unix (`\n`) 换行符差异
- **示例**:
  ```
  标准答案: "hello\nworld"
  用户输出: "hello\r\nworld"
  结果: WA (错误)
  ```

#### 3. **尾随空格**
- **问题**: 用户输出可能有尾随空格，但标准答案没有
- **示例**:
  ```
  标准答案: "42"
  用户输出: "42 "  // 尾随空格
  结果: WA (错误)
  ```

#### 4. **前导空格**
- **问题**: 用户输出可能有前导空格
- **示例**:
  ```
  标准答案: "42"
  用户输出: " 42"  // 前导空格
  结果: WA (错误)
  ```

#### 5. **多个连续空格**
- **问题**: 用户输出可能有多个连续空格，但标准答案只有一个
- **示例**:
  ```
  标准答案: "hello world"
  用户输出: "hello  world"  // 两个空格
  结果: WA (错误)
  ```

## 常见 OJ 系统的处理方式

### 1. **严格模式**（当前实现）
- 完全精确匹配
- 优点: 简单、准确
- 缺点: 可能因为格式问题误判

### 2. **宽松模式**（推荐）
- 去除首尾空白字符
- 统一换行符
- 比较核心内容

### 3. **智能模式**（高级）
- 去除首尾空白
- 统一换行符
- 处理多个连续空格
- 可配置的容错规则

## 建议的改进方案

### 方案 1: 去除首尾空白（推荐）

```go
// 去除首尾空白字符后比较
userOutput := strings.TrimSpace(out.String())
expectedOutput := strings.TrimSpace(tc.Output)

if expectedOutput != userOutput {
    msg = "答案错误"
    WA <- 1
    return
}
```

**优点**:
- 解决尾随空格/换行符问题
- 保持核心内容精确匹配
- 简单易实现

### 方案 2: 统一换行符

```go
// 统一换行符为 \n
normalize := func(s string) string {
    s = strings.ReplaceAll(s, "\r\n", "\n")
    s = strings.ReplaceAll(s, "\r", "\n")
    return strings.TrimSpace(s)
}

userOutput := normalize(out.String())
expectedOutput := normalize(tc.Output)

if expectedOutput != userOutput {
    msg = "答案错误"
    WA <- 1
    return
}
```

### 方案 3: 完整规范化（最宽松）

```go
// 完整的输出规范化
normalizeOutput := func(s string) string {
    // 统一换行符
    s = strings.ReplaceAll(s, "\r\n", "\n")
    s = strings.ReplaceAll(s, "\r", "\n")
    // 去除首尾空白
    s = strings.TrimSpace(s)
    // 将多个连续空白字符替换为单个空格（可选）
    // s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
    return s
}

userOutput := normalizeOutput(out.String())
expectedOutput := normalizeOutput(tc.Output)

if expectedOutput != userOutput {
    msg = "答案错误"
    WA <- 1
    return
}
```

## 当前逻辑总结

**回答你的问题**: 

❌ **不是**"输入和标准答案一样的代码"

✅ **实际逻辑是**:
- **Input**: 输入数据（给用户代码的输入）
- **Output**: 期望输出（标准答案）
- **判断**: 用户代码执行后，将实际输出与期望输出进行**精确字符串比较**

**流程**:
```
用户代码 + Input → 执行 → 实际输出
                              ↓
                        与 Output 比较
                              ↓
                        完全一致 = AC
                        不一致 = WA
```

## 是否需要改进？

### 建议改进的原因:
1. 提高用户体验：避免因格式问题（如尾随空格）导致 WA
2. 符合常见 OJ 系统的做法
3. 减少误判

### 建议的实现:
使用**方案 1**（去除首尾空白），因为：
- 简单有效
- 解决最常见的问题
- 保持核心内容的精确匹配

