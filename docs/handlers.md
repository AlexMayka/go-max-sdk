# Обработчики (Handlers)

Полный справочник по созданию и регистрации обработчиков сообщений.

## Типы обработчиков

### OnStarted()
Обрабатывает событие запуска бота.

```go
func OnStarted(handler Handler)
```

**Пример:**
```go
bot.OnStarted(func(ctx *BotContext) {
    ctx.Reply("🚀 Бот запущен и готов к работе!")
})
```

### OnCommand()
Обрабатывает команды (сообщения начинающиеся с `/`).

```go
func OnCommand(cmd string, handler Handler)
```

**Примеры:**
```go
bot.OnCommand("/start", func(ctx *BotContext) {
    ctx.Reply("Добро пожаловать!")
})

bot.OnCommand("/help", helpHandler)
bot.OnCommand("/settings", settingsHandler)
```

### OnMessage()
Обрабатывает точное соответствие текста сообщения.

```go
func OnMessage(msg string, handler Handler)
```

**Примеры:**
```go
bot.OnMessage("привет", func(ctx *BotContext) {
    ctx.Reply("Привет! Как дела?")
})

bot.OnMessage("пока", func(ctx *BotContext) {
    ctx.Reply("До свидания!")
})
```

### OnPrefix()
Обрабатывает сообщения начинающиеся с префикса.

```go
func OnPrefix(prefix string, handler Handler)
```

**Пример:**
```go
bot.OnPrefix("поиск ", func(ctx *BotContext) {
    query := strings.TrimPrefix(ctx.Text, "поиск ")
    ctx.Reply("Ищу: " + query)
})
```

### OnSuffix()
Обрабатывает сообщения заканчивающиеся суффиксом.

```go
func OnSuffix(suffix string, handler Handler)
```

**Пример:**
```go
bot.OnSuffix(" ?", func(ctx *BotContext) {
    ctx.Reply("Это вопрос! Отвечу позже")
})
```

### OnContains()
Обрабатывает сообщения содержащие подстроку.

```go
func OnContains(sub string, handler Handler)
```

**Пример:**
```go
bot.OnContains("спасибо", func(ctx *BotContext) {
    ctx.Reply("Пожалуйста! 😊")
})
```

### OnCallback()
Обрабатывает callback data от inline кнопок.

```go
func OnCallback(call string, handler Handler)
```

**Пример:**
```go
bot.OnCallback("buy_item", func(ctx *BotContext) {
    ctx.Reply("Покупка подтверждена!")
})
```

### OnRegex()
Обрабатывает сообщения по регулярному выражению.

```go
func OnRegex(regex string, handler Handler)
```

**Примеры:**
```go
// Команда с параметром
bot.OnRegex(`^/calc (\d+) \+ (\d+)$`, calcHandler)

// Hashtag
bot.OnRegex(`^#(\w+)`, hashtagHandler)

// Email
bot.OnRegex(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`, emailHandler)
```

### Any()
Fallback обработчик для всех необработанных сообщений.

```go
func Any(handler Handler)
```

**Пример:**
```go
bot.Any(func(ctx *BotContext) {
    ctx.Reply("❓ Неизвестная команда. Используйте /help")
})
```

### UseState()
Обрабатывает сообщения в определенном FSM состоянии.

```go
func UseState(state string) Router
```

**Пример:**
```go
bot.UseState("waiting_name").Any(func(ctx *BotContext) {
    name := ctx.Text
    ctx.SetData("name", name)
    ctx.SetState("waiting_email")
})
```

## Создание обработчиков

### Простая функция
```go
func startHandler(ctx *BotContext) {
    ctx.Reply("🚀 Бот запущен!")
}

bot.OnCommand("/start", startHandler)
```

### Анонимная функция
```go
bot.OnCommand("/ping", func(ctx *BotContext) {
    ctx.Reply("🏓 Pong!")
})
```

### Метод структуры
```go
type BotHandlers struct {
    userService *UserService
}

func (h *BotHandlers) StartHandler(ctx *BotContext) {
    user := h.userService.GetUser(ctx.UserID)
    ctx.Reply("Привет, " + user.Name + "!")
}

handlers := &BotHandlers{userService: userSvc}
bot.OnCommand("/start", handlers.StartHandler)
```

## Группировка обработчиков

### Создание групп
```go
router := maxsdk.DefaultRouter()

// Публичные команды
router.OnCommand("/start", startHandler)
router.OnCommand("/help", helpHandler)

// Админ группа
adminGroup := router.Group("/admin").Use(adminMiddleware)
adminGroup.OnCommand("/stats", adminStatsHandler)
adminGroup.OnCommand("/users", adminUsersHandler)

// Группа с состоянием
regGroup := router.UseState("registration")
regGroup.Any(registrationHandler)

bot.SetRouter(router)
```

### Middleware для групп
```go
// Логирование для всех обработчиков
router.Use(loggingMiddleware)

// Авторизация только для админ команд
adminGroup := router.Group("/admin").Use(authMiddleware)
```

## Паттерны обработчиков

### Command с параметрами
```go
bot.OnRegex(`^/ban (\d+)$`, func(ctx *BotContext) {
    // Получение параметра из regex группы
    re := regexp.MustCompile(`^/ban (\d+)$`)
    matches := re.FindStringSubmatch(ctx.Text)
    
    if len(matches) > 1 {
        userID := matches[1]
        ctx.Reply("Пользователь " + userID + " заблокирован")
    }
})
```

### Многошаговый диалог
```go
// Шаг 1: Начало регистрации
bot.OnCommand("/register", func(ctx *BotContext) {
    ctx.Reply("Введите ваше имя:")
    ctx.SetState("reg_name")
})

// Шаг 2: Ввод имени
bot.UseState("reg_name").Any(func(ctx *BotContext) {
    name := ctx.Text
    ctx.SetData("name", name)
    ctx.Reply("Введите email:")
    ctx.SetState("reg_email")
})

// Шаг 3: Ввод email
bot.UseState("reg_email").Any(func(ctx *BotContext) {
    email := ctx.Text
    name, _ := ctx.GetData("name")
    
    ctx.Reply("Регистрация завершена!\n👤 " + name + "\n📧 " + email)
    ctx.SetState("") // Сброс состояния
})
```

### Inline клавиатуры
```go
bot.OnCommand("/menu", func(ctx *BotContext) {
    statsPayload := "stats"
    settingsPayload := "settings"
    helpPayload := "help"
    
    keyboard := &models.InlineKeyboard{
        Buttons: [][]models.InlineKeyboardButton{
            {{
                Type: models.InlineKeyboardButtonTypeCallback,
                Text: "📊 Статистика", 
                Payload: &statsPayload,
            }},
            {{
                Type: models.InlineKeyboardButtonTypeCallback,
                Text: "⚙️ Настройки", 
                Payload: &settingsPayload,
            }},
            {{
                Type: models.InlineKeyboardButtonTypeCallback,
                Text: "ℹ️ Помощь", 
                Payload: &helpPayload,
            }},
        },
    }
    ctx.ReplyWithKeyboard("Выберите действие:", keyboard)
})

// Обработка нажатий кнопок
bot.OnCallback("stats", func(ctx *BotContext) {
    ctx.Reply("📊 Ваша статистика:\n• Сообщений: 156\n• Дней активности: 12")
})
```

### Загрузка файлов
```go
bot.Any(func(ctx *BotContext) {
    if ctx.Update.Message != nil && ctx.Update.Message.Body != nil && len(ctx.Update.Message.Body.Attachments) > 0 {
        for _, attachment := range ctx.Update.Message.Body.Attachments {
            switch attachment.Type {
            case "image":
                ctx.Reply("🖼️ Получена картинка")
            case "file":
                ctx.Reply("📎 Получен файл")
            case "video":
                ctx.Reply("🎥 Получено видео")
            }
        }
    }
})
```

## Валидация в обработчиках

### Проверка прав доступа
```go
func adminOnlyHandler(ctx *BotContext) {
    if !isAdmin(ctx.UserID) {
        ctx.Reply("🔒 Недостаточно прав")
        return
    }
    
    ctx.Reply("👑 Админ панель")
}

func isAdmin(userID int64) bool {
    adminIDs := []int64{123456, 789012}
    for _, id := range adminIDs {
        if id == userID {
            return true
        }
    }
    return false
}
```

### Валидация ввода
```go
func emailHandler(ctx *BotContext) {
    email := strings.TrimSpace(ctx.Text)
    
    if !isValidEmail(email) {
        ctx.Reply("❌ Неверный формат email")
        return
    }
    
    ctx.SetData("email", email)
    ctx.Reply("✅ Email сохранен: " + email)
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}
```

### Rate limiting пользователей
```go
var userLastMessage = make(map[int64]time.Time)
var userMessageMutex sync.RWMutex

func rateLimitedHandler(ctx *BotContext) {
    userMessageMutex.Lock()
    lastTime, exists := userLastMessage[ctx.UserID]
    userLastMessage[ctx.UserID] = time.Now()
    userMessageMutex.Unlock()
    
    if exists && time.Since(lastTime) < time.Second {
        ctx.Reply("⏱️ Слишком быстро! Подождите секунду")
        return
    }
    
    ctx.Reply("✅ Сообщение обработано")
}
```

## Обработка ошибок

### Graceful error handling
```go
func safeHandler(ctx *BotContext) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Handler paniced: %v", r)
            ctx.Reply("😅 Произошла ошибка, попробуйте позже")
        }
    }()
    
    // Основная логика
    riskyOperation(ctx)
}
```

### Timeout обработки
```go
func timeoutHandler(ctx *BotContext) {
    timeout := 5 * time.Second
    done := make(chan bool, 1)
    
    go func() {
        // Длительная операция
        time.Sleep(3 * time.Second)
        ctx.Reply("Операция завершена")
        done <- true
    }()
    
    select {
    case <-done:
        // Успешно завершено
    case <-time.After(timeout):
        ctx.Reply("⏱️ Операция заняла слишком много времени")
    }
}
```

## Тестирование обработчиков

### Unit тесты
```go
func TestStartHandler(t *testing.T) {
    // Создаем mock context
    ctx := &BotContext{
        UserID: 12345,
        Text:   "/start",
    }
    
    // Вызываем обработчик
    startHandler(ctx)
    
    // Проверяем результат
    // (требует дополнительной настройки mock'ов)
}
```

### Integration тесты
```go
func TestBotFlow(t *testing.T) {
    bot := maxsdk.NewBot("test-token")
    
    // Регистрируем обработчики
    setupHandlers(bot)
    
    // Симулируем сообщения
    // (требует test transport)
}
```

## Best Practices

### 1. Используйте осмысленные имена
```go
// ❌ Плохо
func h1(ctx *BotContext) {}

// ✅ Хорошо
func startCommandHandler(ctx *BotContext) {}
```

### 2. Группируйте связанные обработчики
```go
// ✅ Хорошо
type UserHandlers struct{}
func (h *UserHandlers) Register(ctx *BotContext) {}
func (h *UserHandlers) Profile(ctx *BotContext) {}
func (h *UserHandlers) Settings(ctx *BotContext) {}
```

### 3. Валидируйте входные данные
```go
func processAmount(ctx *BotContext) {
    amountStr := strings.TrimPrefix(ctx.Text, "/pay ")
    amount, err := strconv.ParseFloat(amountStr, 64)
    
    if err != nil || amount <= 0 {
        ctx.Reply("❌ Неверная сумма")
        return
    }
    
    // Обработка платежа
}
```

### 4. Обрабатывайте ошибки
```go
func sendEmailHandler(ctx *BotContext) {
    err := sendEmail(ctx.UserID)
    if err != nil {
        log.Printf("Email send failed: %v", err)
        ctx.Reply("😅 Не удалось отправить email")
        return
    }
    
    ctx.Reply("✅ Email отправлен")
}
```

### 5. Используйте контекст для отмены
```go
func longOperationHandler(ctx *BotContext) {
    select {
    case result := <-doLongOperation():
        ctx.Reply("Результат: " + result)
    case <-ctx.Ctx.Done():
        // Операция отменена
        return
    }
}
```