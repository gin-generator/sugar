# Sugar - Go Web 框架

基于 Gin 的 Go Web 框架，参考 Laravel 设计模式，提供统一的服务管理和灵活的扩展能力。

## 特性

- 🎯 **服务容器**：统一管理所有服务的生命周期
- 🔌 **服务提供者**：模块化的服务注册和启动机制
- 🎭 **Facade 模式**：提供便捷的静态访问接口
- 🗄️ **多数据库支持**：MySQL、PostgreSQL 等多数据库连接管理
- 💾 **缓存服务**：统一的缓存接口，支持 Redis 等多种驱动
- 📁 **文件存储**：支持本地存储、S3、OSS 等多种存储方式
- 📮 **消息队列**：异步任务处理支持
- 🚀 **多服务类型**：支持 HTTP、WebSocket、gRPC 等多种服务

## 快速开始

### 安装依赖

```bash
go mod download
```

### 配置

编辑 `app/demo/etc/env.yaml` 配置文件：

```yaml
app:
  name: demo
  host: 0.0.0.0
  port: 8888
  env: debug

logger:
  level: debug
  filename: storage/logs/logs.log
  maxSize: 32
  maxBackup: 10
  maxAge: 7
  compress: false

database:
  mysql:
    admin:
      host: 127.0.0.1
      port: 3306
      username: root
      password: your_password
      charset: utf8mb4
      parseTime: true
      multiStatements: true
      loc: Local
```

### 运行

```bash
go run app/demo/demo.go
```

访问 `http://localhost:8888/ping` 测试服务是否正常运行。

## 项目结构

```
.
├── app/                    # 应用层
│   └── demo/              # 示例应用
│       ├── api/           # HTTP API
│       ├── middleware/    # 应用级中间件
│       ├── route/         # 路由
│       └── etc/           # 配置文件
├── bootstrap/             # 启动引导
├── config/                # 配置管理
├── foundation/            # 核心基础（服务容器）
├── providers/             # 服务提供者
├── services/              # 基础服务层
│   ├── database/         # 数据库服务
│   ├── cache/            # 缓存服务
│   ├── storage/          # 文件存储服务
│   ├── queue/            # 消息队列服务
│   └── logger/           # 日志服务
├── middleware/            # 全局中间件
└── model/                 # 数据模型
```

## 使用示例

### 数据库操作

```go
import "github.com/gin-generator/sugar/services/database"

// 使用默认连接
db, _ := database.DB()
var users []User
db.Find(&users)

// 使用指定连接
conn, _ := database.Connection("admin")
conn.Find(&users)
```

### 缓存操作

```go
import "github.com/gin-generator/sugar/services/cache"

ctx := context.Background()

// 设置缓存
cache.Set(ctx, "key", "value", time.Hour)

// 获取缓存
value, _ := cache.Get(ctx, "key")

// 删除缓存
cache.Delete(ctx, "key")
```

### 创建 API

```go
// app/demo/route/route.go
func RegisterApi(e *gin.Engine) {
    e.GET("/users", func(c *gin.Context) {
        db, _ := database.DB()
        var users []User
        db.Find(&users)
        
        c.JSON(http.StatusOK, gin.H{
            "data": users,
        })
    })
}
```

## 文档

- [架构说明](ARCHITECTURE_CN.md) - 详细的架构设计说明
- [使用示例](USAGE_EXAMPLES_CN.md) - 完整的使用示例
- [迁移指南](MIGRATION_GUIDE_CN.md) - 从旧版本迁移指南

## 创建新应用

### 1. 创建应用目录

```bash
mkdir -p app/myapp/{api,middleware,route,etc}
```

### 2. 创建主文件

```go
// app/myapp/myapp.go
package main

import (
    "github.com/gin-generator/sugar/app/myapp/route"
    "github.com/gin-generator/sugar/bootstrap"
    "github.com/gin-generator/sugar/middleware"
)

func main() {
    b := bootstrap.NewBootstrap(
        bootstrap.ServerHttp,
        bootstrap.WithHttpMiddleware(
            middleware.Recovery(),
            middleware.Logger(),
            middleware.Cors(),
        ),
        bootstrap.WithHttpRouter(route.RegisterApi),
    )
    b.Run()
}
```

### 3. 创建配置文件

复制 `app/demo/etc/env.yaml` 到 `app/myapp/etc/env.yaml` 并修改配置。

### 4. 运行应用

```bash
go run app/myapp/myapp.go
```

## 添加自定义服务

### 1. 创建服务

```go
// services/email/manager.go
package email

type Manager struct{}

func NewManager() *Manager {
    return &Manager{}
}

func (m *Manager) Send(to, subject, body string) error {
    // 发送邮件逻辑
    return nil
}
```

### 2. 创建服务提供者

```go
// providers/email.go
package providers

type EmailServiceProvider struct {
    cfg *config.Config
}

func (p *EmailServiceProvider) Register(app *foundation.Application) {
    manager := email.NewManager()
    app.Bind("email", manager)
}

func (p *EmailServiceProvider) Boot(app *foundation.Application) error {
    return nil
}

func (p *EmailServiceProvider) Name() string {
    return "Email"
}
```

### 3. 注册服务提供者

在 `bootstrap/bootstrap.go` 的 `registerProviders` 方法中添加：

```go
b.app.Register(providers.NewEmailServiceProvider(b.cfg))
```

## 许可证

MIT License
