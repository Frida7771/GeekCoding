# GeekCoding 在线判题系统 - 项目功能说明

## 项目概述

GeekCoding 是一个基于 Go + Gin + GORM + MySQL + Redis 的在线判题系统（Online Judge），支持用户提交代码、自动编译、执行和判题。

---

## 核心功能模块

### 1. 用户管理模块 (`service/user.go`)

#### 1.1 用户注册 (`/register` - POST)
- **功能**：用户注册新账号
- **流程**：
  1. 接收邮箱、验证码、用户名、密码、手机号（可选）
  2. 验证邮箱验证码（从 Redis 获取）
  3. 检查邮箱是否已注册
  4. 密码使用 MD5 加密存储
  5. 创建用户记录，默认 `is_admin = 0`（普通用户）
- **返回**：注册成功信息

#### 1.2 用户登录 (`/login` - POST)
- **功能**：用户登录认证
- **流程**：
  1. 接收用户名和密码
  2. 密码 MD5 加密后查询数据库
  3. 生成 JWT Token（包含用户身份、用户名、管理员标识）
  4. Token 有效期 24 小时
- **返回**：JWT Token

#### 1.3 发送验证码 (`/send-code` - POST)
- **功能**：发送邮箱验证码用于注册
- **流程**：
  1. 生成 6 位随机验证码
  2. 将验证码存储到 Redis（有效期 5 分钟）
  3. 通过 SMTP 发送邮件（Gmail）
- **返回**：发送成功信息

#### 1.4 获取用户详情 (`/user-detail` - GET)
- **功能**：根据用户 identity 获取用户信息
- **安全**：自动隐藏密码字段
- **返回**：用户详细信息（不含密码）

#### 1.5 排行榜 (`/rank-list` - GET)
- **功能**：获取用户排行榜
- **排序规则**：
  1. 按完成题目数降序（`finish_problem_num DESC`）
  2. 按提交数升序（`submit_num ASC`）
- **支持分页**：`page` 和 `size` 参数

---

### 2. 题目管理模块 (`service/problem.go`)

#### 2.1 获取题目列表 (`/problem-list` - GET) - 公开接口
- **功能**：分页获取题目列表
- **功能特性**：
  - 支持关键词搜索（标题和内容）
  - 支持按分类筛选（`category_identity`）
  - 预加载分类信息（`ProblemCategorys`）
  - 分页支持
- **返回**：题目列表和总数

#### 2.2 获取题目详情 (`/problem-detail` - GET) - 公开接口
- **功能**：根据题目 identity 获取详细信息
- **包含信息**：
  - 题目基本信息（标题、内容、内存限制、时间限制）
  - 关联的分类信息
  - 测试用例（`TestCases`）
- **返回**：完整的题目信息

#### 2.3 创建题目 (`/admin/problem-create` - POST) - 管理员接口
- **功能**：管理员创建新题目
- **需要参数**：
  - `title`：题目标题
  - `content`：题目描述
  - `max_runtime`：最大运行时间（毫秒）
  - `max_mem`：最大内存限制（MB）
  - `category_ids`：分类 ID 数组（多选）
  - `test_cases`：测试用例数组（JSON 格式，每个包含 `input` 和 `output`）
- **流程**：
  1. 创建 `ProblemBasic` 记录
  2. 创建关联的分类关系（`ProblemCategory`）
  3. 创建测试用例（`TestCase`）
  4. 使用事务确保数据一致性
- **返回**：题目 identity

#### 2.4 更新题目 (`/admin/problem-update` - PUT) - 管理员接口
- **功能**：管理员更新题目信息
- **流程**：
  1. 更新题目基本信息
  2. 删除旧的分类关联，创建新的
  3. 删除旧的测试用例，创建新的
  4. 使用事务确保数据一致性
- **返回**：更新成功信息

---

### 3. 分类管理模块 (`service/category.go`)

#### 3.1 获取分类列表 (`/admin/category-list` - GET) - 管理员接口
- **功能**：分页获取分类列表
- **支持**：关键词搜索（按分类名称）
- **返回**：分类列表和总数

#### 3.2 创建分类 (`/admin/category-create` - POST) - 管理员接口
- **功能**：创建新分类
- **参数**：
  - `name`：分类名称
  - `parent_id`：父分类 ID（可选，支持分类层级）
- **返回**：创建成功信息

#### 3.3 更新分类 (`/admin/category-update` - PUT) - 管理员接口
- **功能**：更新分类信息
- **参数**：
  - `identity`：分类标识（query 参数）
  - `name`：新名称
  - `parent_id`：新父分类 ID
- **返回**：更新成功信息

#### 3.4 删除分类 (`/admin/category-delete` - DELETE) - 管理员接口
- **功能**：删除分类
- **安全检查**：
  - 检查是否有题目使用该分类
  - 如果有题目使用，不允许删除
- **删除方式**：硬删除（`Unscoped().Delete()`）
- **返回**：删除成功信息

---

### 4. 代码提交与判题模块 (`service/submit.go`)

#### 4.1 获取提交列表 (`/submit-list` - GET) - 公开接口
- **功能**：分页获取代码提交记录
- **筛选条件**：
  - `problem_identity`：按题目标识筛选
  - `user_identity`：按用户标识筛选
  - `status`：按提交状态筛选
- **预加载**：关联的题目和用户信息
- **返回**：提交记录列表和总数

#### 4.2 提交代码 (`/user/submit` - POST) - 用户私有接口
- **功能**：用户提交代码进行判题
- **完整流程**：

  **步骤 1：接收和保存代码**
  - 从请求体读取代码
  - 保存到系统临时目录：`/tmp/code/{UUID}/main.go`

  **步骤 2：验证题目**
  - 查询题目是否存在
  - 检查是否有测试用例
  - 预加载测试用例

  **步骤 3：编译代码**
  - 执行 `go build` 编译代码
  - 如果编译失败，直接返回编译错误（状态码 5）
  - 编译成功，生成二进制文件

  **步骤 4：并发执行测试用例**
  - 为每个测试用例启动一个 goroutine
  - 使用 Docker 运行代码：
    - 内存限制：`--memory={max_mem}m`
    - 禁用网络：`--network=none`
    - 超时控制：`--stop-timeout={timeout}s`
  - 输入测试用例的输入数据
  - 捕获输出结果

  **步骤 5：结果判断**
  - **编译错误**（状态码 5）：编译阶段失败
  - **答案错误**（状态码 2）：输出与期望输出不匹配
  - **运行超内存**（状态码 4）：Docker 退出码 137（被 OOM Killer 杀死）
  - **运行超时**（状态码 3）：超过最大运行时间
  - **答案正确**（状态码 1）：所有测试用例通过

  **步骤 6：保存结果**
  - 使用事务保存提交记录
  - 更新用户统计：
    - `submit_num`：提交数 +1
    - `pass_num`：如果答案正确，通过数 +1

- **返回**：提交状态和消息

---

### 5. 认证与授权模块

#### 5.1 用户认证中间件 (`middlewares/auth_user.go`)
- **功能**：验证用户身份
- **流程**：
  1. 从 `Authorization` header 获取 Token
  2. 解析并验证 JWT Token
  3. 将用户信息存储到 context（`c.Set("user", userClaims)`）
  4. 允许访问受保护的用户接口

#### 5.2 管理员认证中间件 (`middlewares/auth_admin.go`)
- **功能**：验证管理员身份
- **流程**：
  1. 从 `Authorization` header 获取 Token
  2. 解析并验证 JWT Token
  3. 检查 `IsAdmin == 1`
  4. 允许访问管理员接口

---

### 6. 辅助功能模块 (`help/helper.go`)

#### 6.1 JWT Token 管理
- **GenerateToken**：生成 JWT Token
  - 包含：用户 identity、用户名、管理员标识
  - 有效期：24 小时
- **AnalyzeToken**：解析和验证 Token

#### 6.2 加密与安全
- **MD5**：MD5 哈希加密（用于密码存储）

#### 6.3 工具函数
- **GetUUID**：生成 UUID（用于唯一标识）
- **GetRandomCode**：生成 6 位随机验证码
- **SendCode**：发送邮件验证码（SMTP）
- **SaveCode**：保存用户提交的代码到临时目录

---

### 7. 数据模型 (`models/`)

#### 7.1 用户模型 (`user_basic.go`)
- **字段**：
  - `identity`：用户唯一标识
  - `username`：用户名
  - `password`：密码（MD5 加密）
  - `email`：邮箱
  - `phone`：手机号
  - `pass_num`：通过题目数
  - `submit_num`：提交总数
  - `is_admin`：是否管理员（0=普通用户，1=管理员）

#### 7.2 题目模型 (`problem_basic.go`)
- **字段**：
  - `identity`：题目唯一标识
  - `title`：题目标题
  - `content`：题目描述
  - `max_runtime`：最大运行时间（毫秒）
  - `max_mem`：最大内存限制（MB）
- **关联**：
  - `TestCases`：测试用例（一对多）
  - `ProblemCategorys`：分类关联（多对多）

#### 7.3 测试用例模型 (`test_case.go`)
- **字段**：
  - `identity`：测试用例唯一标识
  - `problem_identity`：关联的题目标识
  - `input`：输入数据
  - `output`：期望输出

#### 7.4 提交记录模型 (`submit_basic.go`)
- **字段**：
  - `identity`：提交唯一标识
  - `problem_identity`：题目标识
  - `user_identity`：用户标识
  - `path`：代码保存路径
  - `status`：提交状态（1=正确，2=错误，3=超时，4=超内存，5=编译错误）
- **关联**：
  - `ProblemBasic`：关联的题目
  - `User_Basic`：关联的用户

#### 7.5 分类模型 (`category_basic.go`)
- **字段**：
  - `identity`：分类唯一标识
  - `name`：分类名称
  - `parent_id`：父分类 ID（支持分类层级）

#### 7.6 题目分类关联模型 (`problem_category.go`)
- **功能**：多对多关系表
- **字段**：
  - `problem_id`：题目 ID
  - `category_id`：分类 ID

---

## 技术特性

### 1. 代码执行与判题
- **Docker 隔离**：使用 Docker 容器运行用户代码，确保安全隔离
- **内存限制**：通过 Docker 的 `--memory` 参数限制内存使用
- **超时控制**：通过 `--stop-timeout` 和 `time.After` 控制执行时间
- **并发测试**：使用 goroutine 并发执行多个测试用例
- **准确判断**：
  - 编译错误：先编译，失败直接返回
  - 运行时错误：通过 Docker 退出码判断
  - 内存超限：退出码 137（SIGKILL）
  - 答案正确：所有测试用例输出匹配

### 2. 安全性
- **JWT 认证**：使用 JWT Token 进行身份验证
- **密码加密**：MD5 加密存储
- **权限控制**：区分普通用户和管理员
- **代码隔离**：Docker 容器隔离，禁用网络访问

### 3. 数据库设计
- **GORM ORM**：使用 GORM 进行数据库操作
- **自动迁移**：启动时自动创建/更新表结构
- **关联关系**：
  - 题目 ↔ 测试用例（一对多）
  - 题目 ↔ 分类（多对多）
  - 提交 ↔ 题目（多对一）
  - 提交 ↔ 用户（多对一）

### 4. 缓存
- **Redis**：用于存储验证码（5 分钟有效期）

### 5. API 文档
- **Swagger**：完整的 Swagger/OpenAPI 文档
- **自动生成**：使用 `swag init` 生成文档
- **访问地址**：`http://localhost:8080/swagger/index.html`

---

## API 接口总览

### 公开接口（无需认证）
- `GET /problem-list` - 获取题目列表
- `GET /problem-detail` - 获取题目详情
- `GET /user-detail` - 获取用户详情
- `POST /login` - 用户登录
- `POST /send-code` - 发送验证码
- `POST /register` - 用户注册
- `GET /rank-list` - 排行榜
- `GET /submit-list` - 获取提交列表

### 用户私有接口（需要用户认证）
- `POST /user/submit` - 提交代码

### 管理员接口（需要管理员认证）
- `POST /admin/problem-create` - 创建题目
- `PUT /admin/problem-update` - 更新题目
- `GET /admin/category-list` - 获取分类列表
- `POST /admin/category-create` - 创建分类
- `PUT /admin/category-update` - 更新分类
- `DELETE /admin/category-delete` - 删除分类

---

## 提交状态码说明

| 状态码 | 含义 | 说明 |
|--------|------|------|
| 1 | 答案正确 | 所有测试用例通过 |
| 2 | 答案错误 | 输出与期望输出不匹配 |
| 3 | 运行超时 | 超过最大运行时间 |
| 4 | 运行超内存 | 超过最大内存限制 |
| 5 | 编译错误 | 代码编译失败 |

---

## 项目结构

```
GeekCoding/
├── main.go                 # 应用入口
├── router/
│   └── app.go             # 路由配置
├── service/               # 业务逻辑层
│   ├── user.go           # 用户相关服务
│   ├── problem.go        # 题目相关服务
│   ├── submit.go         # 提交与判题服务
│   └── category.go       # 分类管理服务
├── models/                # 数据模型层
│   ├── user_basic.go
│   ├── problem_basic.go
│   ├── test_case.go
│   ├── submit_basic.go
│   ├── category_basic.go
│   ├── problem_category.go
│   └── init.go           # 数据库初始化
├── middlewares/           # 中间件
│   ├── auth_user.go      # 用户认证
│   └── auth_admin.go     # 管理员认证
├── help/                  # 辅助函数
│   └── helper.go         # JWT、加密、工具函数
├── internal/code/         # Docker 相关
│   ├── Dockerfile        # 代码运行器镜像
│   └── docker-runner.sh  # 容器内执行脚本
├── docs/                  # Swagger 文档
└── define/               # 常量定义
```

---

## 数据库表结构

1. **user_basic** - 用户表
2. **problem_basic** - 题目表
3. **test_case** - 测试用例表
4. **submit_basic** - 提交记录表
5. **category_basic** - 分类表
6. **problem_category** - 题目分类关联表

---

## 总结

GeekCoding 是一个功能完整的在线判题系统，实现了：
- ✅ 用户注册、登录、认证
- ✅ 题目管理（CRUD）
- ✅ 分类管理（支持层级）
- ✅ 代码提交与自动判题
- ✅ Docker 隔离执行
- ✅ 内存和时间限制
- ✅ 排行榜系统
- ✅ 完整的 API 文档

系统采用现代化的技术栈，具有良好的安全性和可扩展性。

