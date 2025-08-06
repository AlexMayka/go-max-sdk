# Config Bot

## Описание

Демонстрирует настройку конфигурации бота для оптимизации производительности под конкретные задачи. Показывает как настроить воркеры, тайм-ауты, rate limiting и polling параметры.

## Функциональность

- ⚙️ Настройка количества воркеров
- 📊 Конфигурация очередей сообщений
- ⏱️ Настройка тайм-аутов
- 🚦 Rate limiting для API запросов
- 📡 Параметры long polling
- 📋 Просмотр текущей конфигурации

## Параметры конфигурации

### Обработка сообщений
```go
bot.Config.MaxWorkers = 10         // Количество воркеров для обработки
bot.Config.WorkerQueueSize = 1000  // Размер очереди на воркера
bot.Config.RequestTimeout = 15*time.Second  // Тайм-аут выполнения обработчика
```

### Long polling
```go
bot.Config.PollTimeout = 60    // Серверный тайм-аут polling (секунды)
bot.Config.PollLimit = 50      // Максимум сообщений за запрос
bot.Config.UpdatesBuffer = 100 // Буфер канала обновлений
```

### API клиент
```go
bot.Config.APITimeout = 10*time.Second    // HTTP тайм-аут
bot.Config.APIRateLimit = 30.0            // Запросов в секунду
bot.Config.APIBurstLimit = 100.0          // Максимум запросов в burst
bot.Config.APIMaxRetries = 3              // Повторы при ошибках
bot.Config.APIRetryDelay = 100*time.Millisecond  // Задержка между повторами
```

### Graceful shutdown
```go
bot.Config.ShutdownTimeout = 30*time.Second  // Время на завершение
bot.Config.RetryAttempts = 3                // Повторы обработчиков
bot.Config.RetryDelay = time.Second         // Задержка между повторами
```

## Рекомендации по настройке

### Высоконагруженные боты
```go
bot.Config.MaxWorkers = 20           // Больше воркеров
bot.Config.WorkerQueueSize = 2000    // Большая очередь
bot.Config.APIRateLimit = 50.0       // Выше rate limit
bot.Config.PollLimit = 100           // Больше сообщений за раз
```

### Низконагруженные боты
```go
bot.Config.MaxWorkers = 2            // Меньше воркеров
bot.Config.WorkerQueueSize = 100     // Маленькая очередь
bot.Config.APIRateLimit = 10.0       // Ниже rate limit
bot.Config.RequestTimeout = 5*time.Second  // Короткие тайм-ауты
```

### Медленные обработчики
```go
bot.Config.RequestTimeout = 60*time.Second  // Длинные тайм-ауты
bot.Config.RetryAttempts = 5              // Больше повторов
bot.Config.MaxWorkers = 5                 // Меньше параллельности
```

### Быстрые обработчики  
```go
bot.Config.RequestTimeout = 2*time.Second   // Короткие тайм-ауты
bot.Config.MaxWorkers = 50                 // Высокая параллельность
bot.Config.WorkerQueueSize = 5000          // Большие очереди
```

## Значения по умолчанию

```go
// Воркеры и очереди
MaxWorkers: 4
WorkerQueueSize: 100

// Polling
PollTimeout: 30
PollLimit: 100  
UpdatesBuffer: 100

// API клиент
APITimeout: 5*time.Second
APIRateLimit: 20.0
APIBurstLimit: 40.0
APIMaxRetries: 3
APIRetryDelay: 100*time.Millisecond

// Обработка
RequestTimeout: 10*time.Second
ShutdownTimeout: 30*time.Second
RetryAttempts: 3
RetryDelay: time.Second
```

## Мониторинг производительности

### Команда `/config`
Показывает текущие настройки:
```
Current Config:
MaxWorkers: 10
WorkerQueueSize: 1000
PollLimit: 50
APIRateLimit: 30.0
```

### Логирование
```go
bot.Logging = true  // Включить встроенное логирование
```

## Как тестировать

1. Запустите бота: `go run main.go`
2. Команды:

| Команда | Результат |
|---------|-----------|
| `/start` | Подтверждение работы с кастомными настройками |
| `/config` | Показать текущую конфигурацию |

3. Наблюдайте за логами для анализа производительности

## Типичные сценарии

### Чат-бот с большим количеством пользователей
```go
bot.Config.MaxWorkers = 15
bot.Config.WorkerQueueSize = 2000
bot.Config.APIRateLimit = 40.0
bot.Config.PollLimit = 100
```

### Бот с медленными внешними API
```go
bot.Config.RequestTimeout = 30*time.Second
bot.Config.RetryAttempts = 5
bot.Config.MaxWorkers = 3  // Ограничить параллельность
```

### Бот в dev окружении
```go
bot.Config.APIRateLimit = 5.0      // Меньше нагрузки на API
bot.Config.PollTimeout = 10        // Быстрые циклы
bot.Config.RequestTimeout = 5*time.Second
```

### Production бот
```go
bot.Config.ShutdownTimeout = 60*time.Second  // Больше времени на graceful shutdown
bot.Config.RetryAttempts = 5                // Больше устойчивости
bot.Config.APIMaxRetries = 5                // Больше повторов API
```

## Следующий шаг

Изучите [logger_bot](../logger_bot/) для настройки логирования и мониторинга.