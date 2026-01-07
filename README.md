# Sugar - Go Web Framework

A Go web framework based on Gin, inspired by Laravel design patterns, providing unified service management and flexible extensibility.

## Features

- 🎯 **Service Container**: Unified lifecycle management for all services
- 🔌 **Service Providers**: Modular service registration and bootstrapping mechanism
- 🎭 **Facade Pattern**: Convenient static access interface
- 🗄️ **Multi-Database Support**: Connection management for MySQL, PostgreSQL, and more
- 💾 **Cache Service**: Unified cache interface with support for Redis and other drivers
- 📁 **File Storage**: Support for local storage, S3, OSS, and more
- 📮 **Message Queue**: Async task processing support
- 🚀 **Multiple Service Types**: Support for HTTP, WebSocket, gRPC, and more

## Quick Start

### Install Dependencies

```bash
go mod download
```

### Configuration

Edit the `app/demo/etc/env.yaml` configuration file:

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

### Run

```bash
go run app/demo/demo.go
```

Visit `http://localhost:8888/ping` to test if the service is running properly.

## Project Structure

```
.
├── app/                    # Application layer
│   └── demo/              # Demo application
│       ├── api/           # HTTP API
│       ├── middleware/    # Application-level middleware
│       ├── route/         # Routes
│       └── etc/           # Configuration files
├── bootstrap/             # Bootstrap
├── config/                # Configuration management
├── foundation/            # Core foundation (service container)
├── providers/             # Service providers
├── services/              # Base service layer
│   ├── database/         # Database service
│   ├── cache/            # Cache service
│   ├── storage/          # File storage service
│   ├── queue/            # Message queue service
│   └── logger/           # Logger service
├── middleware/            # Global middleware
└── model/                 # Data models
```

## Usage Examples

### Database Operations

```go
import "github.com/gin-generator/sugar/services/database"

// Use default connection
db := database.DB()
var users []User
db.Find(&users)

// Use specified connection
conn, _ := database.Connection("admin")
conn.Find(&users)
```

### Cache Operations

```go
import "github.com/gin-generator/sugar/services/cache"

ctx := context.Background()

// Set cache
cache.Set(ctx, "key", "value", time.Hour)

// Get cache
value, _ := cache.Get(ctx, "key")

// Delete cache
cache.Delete(ctx, "key")
```

### Create API

```go
// app/demo/route/route.go
func RegisterApi(e *gin.Engine) {
    e.GET("/users", func(c *gin.Context) {
        db := database.DB()
        var users []User
        db.Find(&users)
        
        c.JSON(http.StatusOK, gin.H{
            "data": users,
        })
    })
}
```

## Documentation

- [Architecture](ARCHITECTURE.md) - Detailed architecture design documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Complete usage examples

## Create New Application

### 1. Create Application Directory

```bash
mkdir -p app/myapp/{api,middleware,route,etc}
```

### 2. Create Main File

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

### 3. Create Configuration File

Copy `app/demo/etc/env.yaml` to `app/myapp/etc/env.yaml` and modify the configuration.

### 4. Run Application

```bash
go run app/myapp/myapp.go
```

## Add Custom Service

### 1. Create Service

```go
// services/email/manager.go
package email

type Manager struct{}

func NewManager() *Manager {
    return &Manager{}
}

func (m *Manager) Send(to, subject, body string) error {
    // Email sending logic
    return nil
}
```

### 2. Create Service Provider

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

### 3. Register Service Provider

Add to the `registerProviders` method in `bootstrap/bootstrap.go`:

```go
b.app.Register(providers.NewEmailServiceProvider(b.cfg))
```

## License

MIT License
