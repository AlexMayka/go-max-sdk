# Context Bot

## Описание

Демонстрирует работу с контекстом сообщений, прямую отправку, редактирование и удаление сообщений. Показывает как получать информацию о пользователе и работать с различными методами отправки.

## Функциональность

- 📨 Прямая отправка сообщений по User ID
- ✏️ Редактирование отправленных сообщений
- 🗑️ Удаление сообщений
- 📢 Broadcast рассылка
- 👤 Получение информации о пользователе
- 💬 Различие между User и Chat сообщениями

## Команды

### Отправка сообщений
```bash
/send 12345 Hello there!        # Отправить сообщение пользователю
/broadcast Important news!      # Рассылка группе пользователей  
/demo                          # Демо прямой отправки
```

### Управление сообщениями
```bash
/edit msg_123 New text         # Редактировать сообщение
/delete msg_123                # Удалить сообщение
```

### Информация
```bash
/info                          # Показать контекст сообщения
/start                         # Справка по командам
```

## Работа с контекстом

### Информация о пользователе
```go
ctx.UserID      // ID пользователя
ctx.ChatID      // ID чата (0 для приватных сообщений)
ctx.FirstName   // Имя пользователя
ctx.LastName    // Фамилия (может быть пустой)
ctx.Username    // Username (может быть пустой)
ctx.Text        // Текст сообщения
```

### Отправка сообщений
```go
// Через контекст (ответ отправителю)
ctx.Reply("Response message")

// Прямая отправка конкретному пользователю
bot.SendTextToUser(userID, "Direct message")

// Отправка в чат
bot.SendTextToChat(chatID, "Chat message")
```

### Управление сообщениями
```go
// Редактирование (нужен messageID от API)
ctx.EditMessage(messageID, "Updated text")

// Удаление
ctx.DeleteMessage(messageID)
```

## Регулярные выражения в командах

```go
// Команда с параметрами: /send 12345 Hello
bot.OnRegex(`^/send\s+(\d+)\s+(.+)$`, handler)

// Команда с одним параметром: /delete msg_123
bot.OnRegex(`^/delete\s+(\S+)$`, handler)

// Команда с текстом: /broadcast Some text
bot.OnRegex(`^/broadcast\s+(.+)$`, handler)
```

## Парсинг параметров команд

```go
func sendHandler(ctx *maxsdk.BotContext) {
    // Разбить команду на части
    parts := strings.SplitN(ctx.Text, " ", 3)
    if len(parts) < 3 {
        ctx.Reply("Usage: /send <user_id> <text>")
        return
    }

    // Извлечь параметры
    userID, err := strconv.ParseInt(parts[1], 10, 64)
    text := parts[2]
    
    // Использовать параметры
    bot.SendTextToUser(userID, text)
}
```

## Обработка ошибок

```go
err := bot.SendTextToUser(userID, text)
if err != nil {
    ctx.Reply("Failed to send: " + err.Error())
} else {
    ctx.Reply("Message sent successfully!")
}
```

## Broadcast рассылка

```go
userIDs := []int64{12345, 67890, 11111}
sent := 0

for _, userID := range userIDs {
    err := bot.SendTextToUser(userID, "📢 Broadcast: " + text)
    if err == nil {
        sent++
    }
}

ctx.Reply(fmt.Sprintf("Sent to %d users", sent))
```

## Как тестировать

1. Запустите бота: `go run main.go`
2. Попробуйте команды:

| Команда | Результат |
|---------|-----------|
| `/info` | Показать ваш User ID, Chat ID, имя |
| `/demo` | Отправить демо сообщения |
| `/send [your_user_id] Test` | Получить прямое сообщение |
| `/broadcast News` | Рассылка (проверьте консоль на ошибки) |

## Практическое применение

- **Уведомления**: Отправка уведомлений конкретным пользователям
- **Модерация**: Редактирование/удаление сообщений
- **Рассылки**: Broadcast важной информации
- **Логирование**: Получение контекста для аналитики
- **Персонализация**: Использование имени пользователя

## Следующий шаг

Изучите [config_bot](../config_bot/) для настройки конфигурации бота.