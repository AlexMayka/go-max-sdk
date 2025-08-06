# Router Bot

## Описание

Демонстрирует мощную систему маршрутизации Max Bot SDK с различными типами обработчиков, группировкой маршрутов и middleware.

## Функциональность

- 🎯 Точное совпадение текста
- 🔍 Поиск по префиксу, суффиксу и содержанию
- 📝 Регулярные выражения
- 👥 Группировка маршрутов
- 🔒 Middleware для авторизации
- ⚡ Команды и сообщения

## Типы маршрутизации

### Команды
```go
router.OnCommand("/start", handler)    // /start
router.OnCommand("/help", handler)     // /help
```

### Текстовые сообщения
```go
router.OnMessage("hello", handler)     // точно "hello"
router.OnPrefix("!", handler)          // начинается с "!"
router.OnSuffix("?", handler)          // заканчивается на "?"
router.OnContains("bot", handler)      // содержит "bot"
router.OnRegex(`^\d+$`, handler)       // только цифры
```

### Группы маршрутов
```go
adminGroup := router.Group("/admin").Use(authMiddleware)
adminGroup.OnCommand("/status", handler)  // /admin/status с проверкой доступа

publicGroup := router.Group("/public")
publicGroup.OnCommand("/help", handler)   // /public/help без ограничений
```

## Middleware

```go
func authMiddleware(next maxsdk.Handler) maxsdk.Handler {
    return func(ctx *maxsdk.BotContext) {
        if adminUser, _ := ctx.GetData("admin"); adminUser != "true" {
            ctx.Reply("Access denied. Use /login first")
            return
        }
        next(ctx)
    }
}
```

## Как тестировать

1. Запустите бота: `go run main.go`
2. Попробуйте разные сообщения:

| Сообщение | Обработчик | Результат |
|-----------|------------|-----------|
| `hello` | OnMessage | "Exact match: hello!" |
| `!test` | OnPrefix | "Message starts with '!'" |
| `test?` | OnSuffix | "Question detected" |
| `my bot rocks` | OnContains | "You mentioned 'bot'" |
| `123` | OnRegex | "You sent a number" |
| `/public/help` | Group + Command | Список методов роутера |
| `/admin/status` | Protected | "Access denied" (без авторизации) |

## Что изучить

- **Приоритет маршрутов**: более специфичные срабатывают первыми
- **Группировка**: логическая организация команд
- **Middleware**: сквозная функциональность (авторизация, логирование)
- **Кастомный роутер**: `maxsdk.DefaultRouter()` + `bot.SetRouter()`

## Следующий шаг

Изучите [fsm_bot](../fsm_bot/) для работы с состояниями пользователей.