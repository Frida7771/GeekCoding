# Swagger 文档修复指南

## 问题原因

Swagger 文档显示为空的主要原因是：
1. **Swagger 注释格式不正确**：缺少空格或格式错误
2. **缺少通用信息**：main.go 中没有 Swagger 基本信息
3. **部分函数缺少注释**：`GetProblemDetail` 和 `Login` 函数没有 Swagger 注释
4. **文档未重新生成**：修改注释后需要重新生成文档

## 已修复的问题

### 1. 修复了 main.go
添加了 Swagger 通用信息注释：
```go
// @title           GeekCoding API
// @version         1.0
// @description     This is a GeekCoding Online Judge API server.
// @host      localhost:8080
// @BasePath  /
```

### 2. 修复了所有 service 文件的注释格式
- 修正了 `@Tags` 前的空格问题
- 添加了 `@Description`、`@Accept`、`@Produce` 标签
- 为所有函数添加了完整的 Swagger 注释
- 修正了 `@Success` 和 `@Failure` 的格式

## 重新生成 Swagger 文档

### 步骤 1: 安装 swag 工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

如果 `$GOPATH/bin` 不在 PATH 中，需要添加：
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 步骤 2: 重新生成文档

在项目根目录执行：

```bash
swag init
```

或者指定输出目录：

```bash
swag init -g main.go -o ./docs
```

### 步骤 3: 验证生成结果

检查 `docs/swagger.json` 和 `docs/swagger.yaml` 文件，应该包含 `paths` 信息。

### 步骤 4: 运行项目并访问 Swagger

```bash
go run main.go
```

然后访问：`http://localhost:8080/swagger/index.html`

## 常见问题

### 问题 1: swag: command not found

**解决**：
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 确保 GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin

# 验证安装
swag version
```

### 问题 2: 生成后 paths 仍然为空

**可能原因**：
1. 注释格式仍然有问题
2. 路由路径不匹配

**检查**：
- 确保所有注释都以 `//` 开头，且 `@` 前有空格
- 确保 `@Router` 中的路径与 `router/app.go` 中的路径一致

### 问题 3: Swagger UI 显示 404

**解决**：
- 确保 `router/app.go` 中有正确的路由：
  ```go
  r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
  ```
- 确保导入了 docs 包：
  ```go
  _ "GeekCoding/docs"
  ```

## 验证修复

运行以下命令验证：

```bash
# 1. 检查 swag 是否安装
swag version

# 2. 重新生成文档
swag init

# 3. 检查生成的文档
cat docs/swagger.json | grep -A 5 "paths"

# 4. 运行项目
go run main.go

# 5. 访问 Swagger UI
# 浏览器打开: http://localhost:8080/swagger/index.html
```

## 修复后的注释格式示例

```go
// GetProblemList
// @Tags         Public Method
// @Summary      Get Problem List
// @Description  Get a list of problems with pagination
// @Accept       json
// @Produce      json
// @Param        page  query     int     false  "page number"
// @Param        size  query     int     false  "page size"
// @Success      200   {object}    map[string]interface{}
// @Router       /problem-list [get]
func GetProblemList(c *gin.Context) {
    // ...
}
```

## 注意事项

1. **注释位置**：Swagger 注释必须紧贴在函数定义之前
2. **空格要求**：`//` 和 `@` 之间必须有空格
3. **路径匹配**：`@Router` 中的路径必须与 `router/app.go` 中的路径完全一致
4. **重新生成**：每次修改注释后都需要运行 `swag init` 重新生成文档

