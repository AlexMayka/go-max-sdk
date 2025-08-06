# Конфигурация

Полный справочник по настройке Max Bot SDK.

## Структура Config

```go
type Config struct {
    // Engine настройки
    MaxWorkers      int           
    WorkerQueueSize int           
    RequestTimeout  time.Duration 
    ShutdownTimeout time.Duration 
    RetryAttempts   int           
    RetryDelay      time.Duration 
    
    // Transport настройки
    PollTimeout    int           
    PollLimit      int           
    UpdatesBuffer  int           
    PollRetryDelay time.Duration 
    PollMaxRetries int           
    Marker         int64
    
    // API Client настройки
    APIRateLimit  float64       
    APIBurstLimit float64       
    APITimeout    time.Duration 
    APIMaxRetries int           
    APIRetryDelay time.Duration 
    
    // Webhook настройки (будущее)
    WebhookURL    string 
    WebhookSecret string 
}
```

## Engine настройки

Управляют производительностью и обработкой сообщений.

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|-------------|----------|
| `MaxWorkers` | `int` | 20 | Количество параллельных воркеров для обработки сообщений |
| `WorkerQueueSize` | `int` | 100 | Размер буфера между транспортом и воркерами |
| `RequestTimeout` | `time.Duration` | 30s | Максимальное время выполнения одного хендлера |
| `ShutdownTimeout` | `time.Duration` | 10s | Время ожидания завершения воркеров при остановке |
| `RetryAttempts` | `int` | 3 | Количество попыток выполнить хендлер при ошибке |
| `RetryDelay` | `time.Duration` | 1s | Задержка между попытками retry |

### Рекомендации по Engine

**Низкая нагрузка (< 100 сообщений/мин):**
```go
bot.Config.MaxWorkers = 10
bot.Config.WorkerQueueSize = 50
```

**Средняя нагрузка (100-1000 сообщений/мин):**
```go
bot.Config.MaxWorkers = 20          // По умолчанию
bot.Config.WorkerQueueSize = 100    // По умолчанию
```

**Высокая нагрузка (> 1000 сообщений/мин):**
```go
bot.Config.MaxWorkers = 50
bot.Config.WorkerQueueSize = 500
bot.Config.RequestTimeout = 15 * time.Second  // Сокращаем с 30s до 15s для высокой нагрузки
```

## Transport настройки

Управляют получением обновлений от сервера.

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|-------------|----------|
| `PollTimeout` | `int` | 30 | Timeout long polling запросов (секунды) |
| `PollLimit` | `int` | 100 | Максимум updates за один запрос |
| `UpdatesBuffer` | `int` | 1000 | Размер буфера канала updates |
| `PollRetryDelay` | `time.Duration` | 1s | Задержка между попытками polling при ошибках |
| `PollMaxRetries` | `int` | 5 | Максимум retry попыток для polling |
| `Marker` | `int64` | 0 | Стартовый маркер для polling |

### Оптимизация Transport

**Стабильное соединение:**
```go
bot.Config.PollTimeout = 60      // Длинные запросы
bot.Config.PollLimit = 100       // Стандартный лимит
```

**Нестабильное соединение:**
```go
bot.Config.PollTimeout = 10      // Короткие запросы
bot.Config.PollRetryDelay = 5 * time.Second
bot.Config.PollMaxRetries = 10
```

**Высокая нагрузка:**
```go
bot.Config.PollLimit = 500       // Больше updates за запрос
bot.Config.UpdatesBuffer = 5000  // Больший буфер
```

## API Client настройки

Управляют rate limiting для API запросов.

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|-------------|----------|
| `APIRateLimit` | `float64` | 20.0 | Количество запросов в секунду |
| `APIBurstLimit` | `float64` | 20.0 | Максимум токенов в bucket |
| `APITimeout` | `time.Duration` | 30s | Timeout для HTTP запросов |
| `APIMaxRetries` | `int` | 3 | Максимум retry попыток для API запросов |
| `APIRetryDelay` | `time.Duration` | 100ms | Базовая задержка между retry попытками |

### Rate Limiting стратегии

**Консервативный подход:**
```go
bot.Config.APIRateLimit = 10.0   // 10 req/sec
bot.Config.APIBurstLimit = 10.0
```

**Агрессивный подход:**
```go
bot.Config.APIRateLimit = 30.0   // 30 req/sec
bot.Config.APIBurstLimit = 50.0  // Больше burst
```

**Production настройки:**
```go
bot.Config.APIRateLimit = 20.0
bot.Config.APITimeout = 60 * time.Second
bot.Config.APIMaxRetries = 5
```

## Примеры конфигураций

### Development
```go
bot := maxsdk.NewBot("token")
bot.Config.MaxWorkers = 3
bot.Config.APIRateLimit = 5.0
bot.Config.RequestTimeout = 5 * time.Second
bot.Logging = true
```

### Production
```go
bot := maxsdk.NewBot("token")
bot.Config.MaxWorkers = 50
bot.Config.WorkerQueueSize = 1000
bot.Config.APIRateLimit = 25.0
bot.Config.APITimeout = 60 * time.Second
bot.Config.ShutdownTimeout = 15 * time.Second  // Увеличиваем с 10s до 15s
bot.Logging = false  // Или кастомный логгер
```

### High Load
```go
bot := maxsdk.NewBot("token")
bot.Config.MaxWorkers = 100
bot.Config.WorkerQueueSize = 2000
bot.Config.APIRateLimit = 30.0
bot.Config.APIBurstLimit = 100.0
bot.Config.PollLimit = 500
bot.Config.UpdatesBuffer = 10000
```

## Мониторинг производительности

### Метрики для отслеживания

1. **Queue Size** - размер очереди воркеров
2. **Processing Time** - время обработки сообщений
3. **API Rate** - частота API запросов
4. **Error Rate** - частота ошибок retry

### Индикаторы проблем

**Переполнение очереди:**
```
Увеличить MaxWorkers или WorkerQueueSize
```

**Медленная обработка:**
```
Уменьшить RequestTimeout, оптимизировать хендлеры
```

**API ошибки:**
```
Уменьшить APIRateLimit, увеличить APIRetryDelay
```

## Динамическая конфигурация

```go
// Изменение во время работы
bot.Config.MaxWorkers = 30

// Применение изменений (требует перезапуска engine)
bot.Stop()
bot.Start()
```

## Validation конфигурации

SDK автоматически валидирует конфигурацию:

- `MaxWorkers` >= 1
- `WorkerQueueSize` >= 1  
- `APIRateLimit` > 0
- `APITimeout` > 0
- `PollTimeout` >= 1

При неверных значениях используются значения по умолчанию.