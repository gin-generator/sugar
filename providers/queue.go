package providers

import (
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/queue"
)

// QueueServiceProvider queue service provider
type QueueServiceProvider struct{}

// NewQueueServiceProvider creates a queue service provider
func NewQueueServiceProvider() *QueueServiceProvider {
	return &QueueServiceProvider{}
}

// Register registers the service
func (p *QueueServiceProvider) Register(app *foundation.Application) {
	manager := queue.NewManager()
	app.Bind(ServiceQueue, manager)
}

// Boot boots the service
func (p *QueueServiceProvider) Boot(app *foundation.Application) error {
	_ = foundation.MustMake[*queue.Manager](app, ServiceQueue)

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
