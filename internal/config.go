package internal

import "time"

// Config содержит все настройки для работы бота
type Config struct {
	// === ENGINE НАСТРОЙКИ ===
	// Управляют производительностью и обработкой сообщений

	MaxWorkers      int           // Количество параллельных воркеров для обработки сообщений (больше = быстрее, но больше нагрузка)
	WorkerQueueSize int           // Размер буфера между транспортом и воркерами (больше = меньше потерь при пиках)
	RequestTimeout  time.Duration // Максимальное время выполнения одного хендлера (защита от зависших хендлеров)
	ShutdownTimeout time.Duration // Время ожидания завершения воркеров при остановке (graceful shutdown)
	RetryAttempts   int           // Количество попыток выполнить хендлер при ошибке (защита от временных сбоев)
	RetryDelay      time.Duration // Задержка между попытками retry (защита от спама retry)

	// === TRANSPORT НАСТРОЙКИ ===
	// Управляют получением обновлений от сервера

	PollTimeout    int           // Timeout long polling запросов в секундах (30-60 оптимально, больше = меньше запросов)
	PollLimit      int           // Максимум updates за один запрос (100 стандарт, больше = эффективнее при высокой нагрузке)
	UpdatesBuffer  int           // Размер буфера канала updates (больше = защита от потерь при пиках)
	PollRetryDelay time.Duration // Задержка между попытками polling при ошибках (защита от спама API)
	PollMaxRetries int           // Максимум retry попыток для polling (защита от бесконечных ошибок)
	Marker         int64

	// === API CLIENT НАСТРОЙКИ ===
	// Управляют rate limiting для API запросов

	APIRateLimit  float64       // Количество запросов в секунду (20 = 20 req/sec)
	APIBurstLimit float64       // Максимум токенов в bucket (больше = больше burst запросов)
	APITimeout    time.Duration // Timeout для HTTP запросов
	APIMaxRetries int           // Максимум retry попыток для API запросов (защита от временных API ошибок)
	APIRetryDelay time.Duration // Базовая задержка между retry попытками API запросов

	// === WEBHOOK НАСТРОЙКИ (для будущего) ===
	WebhookURL    string // URL для получения webhook
	WebhookSecret string // Secret для валидации webhook
}

// DefaultConfig возвращает оптимальную конфигурацию для большинства ботов
func DefaultConfig() *Config {
	return &Config{
		// Engine - настроено для средней нагрузки
		MaxWorkers:      20,               // 20 воркеров достаточно для большинства ботов
		WorkerQueueSize: 100,              // Буфер на 100 сообщений
		RequestTimeout:  30 * time.Second, // 30 сек на обработку одного сообщения
		ShutdownTimeout: 10 * time.Second, // 10 сек на graceful shutdown
		RetryAttempts:   3,                // 3 попытки при ошибке
		RetryDelay:      1 * time.Second,  // 1 сек между попытками

		// Transport - оптимизировано для стабильности
		PollTimeout:    30,              // 30 сек timeout
		PollLimit:      100,             // 100 updates за запрос
		UpdatesBuffer:  1000,            // Буфер на 1000 updates
		PollRetryDelay: 1 * time.Second, // 1 сек между retry при ошибках
		PollMaxRetries: 5,               // 5 попыток при ошибках API

		// API Client - консервативные лимиты
		APIRateLimit:  20.0,                   // 20 запросов в секунду
		APIBurstLimit: 20.0,                   // Burst до 20 запросов
		APITimeout:    30 * time.Second,       // 30 сек timeout для HTTP
		APIMaxRetries: 3,                      // 3 retry попытки для API запросов
		APIRetryDelay: 100 * time.Millisecond, // 100ms базовая задержка между retry
	}
}
