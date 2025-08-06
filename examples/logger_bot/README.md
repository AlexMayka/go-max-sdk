# Logger Bot

## Описание

Демонстрирует работу с системой логирования Max Bot SDK. Показывает как использовать встроенный логгер, управлять уровнями логирования и интегрировать собственную систему логирования.

## Функциональность

- 📝 Логирование на разных уровнях (Debug, Info, Warn, Error)
- 🔧 Управление включением/отключением логирования
- 📊 Структурированное логирование с компонентами и событиями
- 🖥️ Консольный вывод логов
- 📈 Логирование действий пользователей

## Команды

| Команда | Описание |
|---------|----------|
| `/start` | Запуск бота с логированием |
| `/log` | Тест всех уровней логирования |
| `/stats` | Запрос статистики с логированием |
| `/disable` | Отключить логирование |
| `/enable` | Включить логирование |

## Уровни логирования

### Debug
```go
bot.Logger.Debug("component", "event", "details")
```
- Детальная отладочная информация
- Обычно отключена в production

### Info  
```go
bot.Logger.Info("component", "event", "details")
```
- Информационные сообщения
- Важные события в работе бота

### Warn
```go
bot.Logger.Warn("component", "event", "details")
```
- Предупреждения о потенциальных проблемах
- Не критичные ошибки

### Error
```go
bot.Logger.Error("component", "event", "details")
```
- Критические ошибки
- Проблемы требующие внимания

## Структура логов

Каждый лог содержит три компонента:

```go
bot.Logger.Info("bot", "start_command", "User started the bot")
//             ^^^    ^^^^^^^^^^^^^^   ^^^^^^^^^^^^^^^^^^^^
//          component     event           details
```

### Компоненты системы
- `"bot"` - основная логика бота
- `"api_client"` - API запросы  
- `"transport"` - long polling транспорт
- `"registry"` - маршрутизация
- `"engine"` - обработка сообщений

### Типы событий
- `"start_command"` - пользовательские команды
- `"message_sent"` - отправка сообщений
- `"request_start"` - начало API запроса
- `"rate_limit_wait"` - ожидание rate limit
- `"retry_attempt"` - повтор запроса

## Управление логированием

### Включение/отключение
```go
bot.Logging = true   // Включить встроенное логирование
bot.Logging = false  // Отключить встроенное логирование
```

### Проверка состояния
```go
if bot.Logging {
    bot.Logger.Info("bot", "logging_active", "Logging is enabled")
}
```

## Кастомный логгер

Вы можете заменить встроенный логгер:

```go
type CustomLogger struct{}

func (l *CustomLogger) Debug(component, event, details string) {
    // Ваша логика для debug
}

func (l *CustomLogger) Info(component, event, details string) {
    // Ваша логика для info
}

func (l *CustomLogger) Warn(component, event, details string) {
    // Ваша логика для warn
}

func (l *CustomLogger) Error(component, event, details string) {
    // Ваша логика для error
}

// Установка кастомного логгера
bot.Logger = &CustomLogger{}
```

## Интеграция с внешними системами

### Logrus
```go
import "github.com/sirupsen/logrus"

type LogrusLogger struct {
    logger *logrus.Logger
}

func (l *LogrusLogger) Info(component, event, details string) {
    l.logger.WithFields(logrus.Fields{
        "component": component,
        "event":     event,
    }).Info(details)
}
```

### Zap
```go
import "go.uber.org/zap"

type ZapLogger struct {
    logger *zap.Logger
}

func (l *ZapLogger) Info(component, event, details string) {
    l.logger.Info(details,
        zap.String("component", component),
        zap.String("event", event),
    )
}
```

### Slog (Go 1.21+)
```go
import "log/slog"

type SlogLogger struct {
    logger *slog.Logger
}

func (l *SlogLogger) Info(component, event, details string) {
    l.logger.Info(details,
        "component", component,
        "event", event,
    )
}
```

## Автоматическое логирование SDK

SDK автоматически логирует:

- 🌐 **API запросы**: `"api_client", "request_start"`
- 🔄 **Retry попытки**: `"api_client", "retry_attempt"`  
- 🚦 **Rate limiting**: `"api_client", "rate_limit_wait"`
- 📡 **Transport**: `"transport", "polling_start"`
- ⚙️ **Engine**: `"engine", "worker_start"`

## Пример вывода

```
2024/01/15 10:30:15 [INFO] bot.start_command: User started the bot
2024/01/15 10:30:20 [DEBUG] bot.debug_test: Debug message from user  
2024/01/15 10:30:20 [INFO] bot.info_test: Info message from user
2024/01/15 10:30:20 [WARN] bot.warn_test: Warning message from user
2024/01/15 10:30:20 [ERROR] bot.error_test: Error message from user
2024/01/15 10:30:25 [INFO] bot.stats_requested: Stats requested by user 12345 (John)
```

## Как тестировать

1. Запустите бота: `go run main.go`
2. Наблюдайте за консолью - там появятся логи
3. Попробуйте команды:

| Команда | Что смотреть в консоли |
|---------|------------------------|
| `/start` | `[INFO] bot.start_command` |
| `/log` | Все 4 уровня логирования |
| `/stats` | `[INFO] bot.stats_requested` с данными пользователя |
| `/disable` | Логи перестанут появляться |
| `/enable` | Логи возобновятся |

## Best Practices

1. **Структурированность**: Используйте консистентные компоненты и события
2. **Детализация**: В details включайте полезную информацию
3. **Уровни**: Debug для разработки, Info для production
4. **Performance**: В high-load окружении используйте асинхронные логгеры
5. **Безопасность**: Не логируйте токены и персональные данные

## Следующий шаг

Изучите [advanced_bot](../advanced_bot/) для продвинутых техник разработки ботов.