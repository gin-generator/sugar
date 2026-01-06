package providers

import (
	"fmt"
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/queue"
)

// QueueServiceProvider queue service provider
type QueueServiceProvider struct {
	cfg *config.Config
}

// NewQueueServiceProvider creates a queue service provider
func NewQueueServiceProvider(cfg *config.Config) *QueueServiceProvider {
	return &QueueServiceProvider{cfg: cfg}
}

// Register registers the service
func (p *QueueServiceProvider) Register(app *foundation.Application) {
	manager := queue.NewManager()
	app.Bind("queue", manager)
}

// Boot boots the service
func (p *QueueServiceProvider) Boot(app *foundation.Application) error {
	service, ok := app.Make("queue")
	if !ok {
		return fmt.Errorf("queue service not found")
	}

	_ = service.(*queue.Manager)

	// Initialize queue connections here, such as Redis, RabbitMQ, etc.
	// Example:
	// redisQueue := queue.NewRedisQueue(...)
	// manager.AddConnection("redis", redisQueue)

	return nil
}

// Name returns the service provider name
func (p *QueueServiceProvider) Name() string {
	return "Queue"
}
