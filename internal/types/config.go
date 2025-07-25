package types

import "time"

type EngineConfig struct {
	// Worker Pool настройки
	MaxWorkers      int // Максимальное количество воркеров для обработки сообщений
	WorkerQueueSize int // Размер очереди для каждого воркера

	// Таймауты
	RequestTimeout  time.Duration // Максимальное время обработки одного сообщения
	ShutdownTimeout time.Duration // Время ожидания graceful shutdown

	// Retry политика
	RetryAttempts int           // Количество попыток при ошибке обработки
	RetryDelay    time.Duration // Задержка между попытками

	// Rate limiting
	RateLimit  int           // Максимум сообщений в секунду (0 = без лимита)
	BurstLimit int           // Размер burst для rate limiter
	RatePeriod time.Duration // Период для rate limiting

	// Monitoring и логирование
	EnableMetrics     bool // Включить сбор метрик производительности
	EnableHealthCheck bool // Включить health check endpoint
	LogErrors         bool // Логировать ошибки обработки
}

// DefaultEngineConfig возвращает конфигурацию по умолчанию
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		MaxWorkers:      100,
		WorkerQueueSize: 100,

		RequestTimeout:  30 * time.Second,
		ShutdownTimeout: 10 * time.Second,

		RetryAttempts: 3,
		RetryDelay:    1 * time.Second,

		RateLimit:  0,
		BurstLimit: 10,
		RatePeriod: time.Second,

		EnableMetrics:     false,
		EnableHealthCheck: false,
		LogErrors:         true,
	}
}

// WorkerStats содержит статистику работы воркеров
type WorkerStats struct {
	ActiveWorkers   int           // Количество активных воркеров
	QueuedMessages  int           // Сообщений в очереди
	ProcessedTotal  uint64        // Всего обработано сообщений
	ErrorsTotal     uint64        // Всего ошибок
	AverageLatency  time.Duration // Средняя задержка обработки
	LastProcessedAt time.Time     // Время последней обработки
}

// HealthStatus статус здоровья engine
type HealthStatus struct {
	Status    string        `json:"status"`    // "healthy", "degraded", "unhealthy"
	Uptime    time.Duration `json:"uptime"`    // Время работы
	Stats     WorkerStats   `json:"stats"`     // Статистика воркеров
	Errors    []string      `json:"errors"`    // Последние ошибки
	Timestamp time.Time     `json:"timestamp"` // Время проверки
}
