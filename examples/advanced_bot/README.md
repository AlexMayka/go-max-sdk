# Advanced Bot

## Описание

Комплексный пример, демонстрирующий продвинутые возможности Max Bot SDK: многоуровневую архитектуру, системы регистрации и авторизации, группировку маршрутов, middleware и обработку сложных сценариев.

## Функциональность

- 🏗️ **Модульная архитектура** с разделением на компоненты
- 👤 **Система регистрации** с валидацией данных
- 🔐 **Admin панель** с авторизацией
- 🛡️ **Middleware** для логирования и rate limiting
- 📝 **Обработка файлов** через regex
- 🏷️ **Система hashtag** для трекинга
- ⚙️ **Кастомная конфигурация** для production

## Архитектура

```
Advanced Bot
├── Public Routes
│   ├── /start
│   └── /help
├── Registration System
│   ├── /register
│   ├── reg_name (state)
│   ├── reg_email (state)
│   └── /profile
├── Admin Panel
│   ├── /admin/login <pass>
│   ├── /admin/stats (protected)
│   └── /admin/users (protected)
├── Special Handlers
│   ├── /process <text>
│   └── #hashtag
└── Middleware Stack
    ├── Logging
    ├── Rate Limiting
    └── Admin Auth
```

## Компоненты системы

### 1. Система регистрации

**Многошаговая регистрация с валидацией:**

```go
/register → reg_name → reg_email → completed
```

**Валидация данных:**
- Имя: минимум 2 символа, только буквы
- Email: полная проверка формата

**FSM состояния:**
- `reg_name` - ввод имени
- `reg_email` - ввод email
- Данные сохраняются между шагами

### 2. Admin панель

**Двухуровневая авторизация:**
```go
/admin/login secret123  // Получить доступ
/admin/stats           // Только для админов
/admin/users           // Только для админов
```

**Middleware защита:**
- Проверка пароля
- Сохранение статуса авторизации
- Блокировка неавторизованных запросов

### 3. Обработчики с Regex

**Динамические команды:**
```go
/process Hello World    // → "HELLO WORLD"
#programming           // → "Hashtag tracked: #programming"
```

## Middleware система

### Logging Middleware
```go
func loggingMiddleware(next maxsdk.Handler) maxsdk.Handler {
    return func(ctx *maxsdk.BotContext) {
        // Логирование до обработки
        start := time.Now()
        next(ctx)
        // Логирование после обработки
        duration := time.Since(start)
        log.Printf("Processed in %v", duration)
    }
}
```

### Rate Limiting Middleware
```go
func rateLimitMiddleware(next maxsdk.Handler) maxsdk.Handler {
    return func(ctx *maxsdk.BotContext) {
        // Проверка лимитов
        if checkRateLimit(ctx.UserID) {
            next(ctx)
        } else {
            ctx.Reply("Rate limit exceeded")
        }
    }
}
```

### Admin Auth Middleware
```go
func adminAuthMiddleware(next maxsdk.Handler) maxsdk.Handler {
    return func(ctx *maxsdk.BotContext) {
        // Обработка логина
        if strings.HasPrefix(ctx.Text, "/admin/login ") {
            handleLogin(ctx)
            return
        }
        
        // Проверка доступа
        if !isAdmin(ctx) {
            ctx.Reply("Access denied")
            return
        }
        
        next(ctx)
    }
}
```

## Пошаговые сценарии

### Регистрация пользователя

1. **Начало**: `/register`
   ```
   📝 Registration started!
   Step 1/2: Enter your full name:
   ```

2. **Ввод имени**: `John Doe`
   ```
   ✅ Name: John Doe
   Step 2/2: Enter your email:
   ```

3. **Ввод email**: `john@example.com`
   ```
   🎉 Registration complete!
   👤 John Doe
   📧 john@example.com
   ```

4. **Просмотр профиля**: `/profile`
   ```
   📋 Your Profile:
   👤 John Doe
   📧 john@example.com
   ```

### Доступ к admin панели

1. **Попытка без авторизации**: `/admin/stats`
   ```
   🔒 Admin access required. Use '/admin/login <password>'
   ```

2. **Неправильный пароль**: `/admin/login wrong`
   ```
   ❌ Invalid password
   ```

3. **Правильный пароль**: `/admin/login secret123`
   ```
   ✅ Admin access granted
   ```

4. **Доступ к статистике**: `/admin/stats`
   ```
   📊 Admin Stats:
   • Users: 156
   • Messages: 2,341
   • Uptime: 5d 12h
   ```

## Конфигурация production

```go
bot.Config.MaxWorkers = 15              // Больше воркеров
bot.Config.WorkerQueueSize = 500        // Большая очередь
bot.Config.RequestTimeout = 20*time.Second  // Длинные тайм-ауты
bot.Config.APIRateLimit = 25.0          // Оптимальный rate limit
```

## Валидация и безопасность

### Валидация имени
```go
if len(name) < 2 {
    return "Name too short"
}
if matched, _ := regexp.MatchString(`^[a-zA-Z\s]+$`, name); !matched {
    return "Letters only"
}
```

### Валидация email
```go
emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
if matched, _ := regexp.MatchString(emailRegex, email); !matched {
    return "Invalid email format"
}
```

### Безопасность admin панели
```go
// Простая авторизация (в production используйте JWT/sessions)
if password == "secret123" {
    ctx.SetData("admin", "true")  // Сохранить статус
}

// Проверка доступа
if admin, _ := ctx.GetData("admin"); admin != "true" {
    return "Access denied"
}
```

## Как тестировать

1. **Запуск**: `go run main.go`

2. **Базовые команды**:
   - `/start` - приветствие и список команд
   - `/help` - подробная справка

3. **Регистрация**:
   - `/register` → введите имя → введите email
   - `/profile` - просмотр данных

4. **Admin панель**:
   - `/admin/stats` (будет отказ)
   - `/admin/login secret123` (получить доступ)
   - `/admin/stats` (теперь работает)

5. **Специальные функции**:
   - `/process Hello World` → обработка текста
   - `#programming` → трекинг hashtag

## Расширения

### Добавление новых middleware
```go
router.Use(corsMiddleware, authMiddleware, metricsMiddleware)
```

### Новые группы маршрутов
```go
apiGroup := router.Group("/api").Use(apiAuthMiddleware)
webhookGroup := router.Group("/webhook").Use(webhookMiddleware)
```

### Интеграция с базой данных
```go
func registerEmailHandler(ctx *maxsdk.BotContext) {
    email := ctx.Text
    
    // Сохранить в БД
    err := db.SaveUser(ctx.UserID, name, email)
    if err != nil {
        ctx.Reply("Registration failed")
        return
    }
    
    ctx.Reply("Registration successful!")
}
```

## Best Practices

1. **Модульность**: Разделяйте логику на функции
2. **Валидация**: Всегда проверяйте пользовательский ввод  
3. **Безопасность**: Не храните пароли в коде
4. **Ошибки**: Обрабатывайте все возможные сценарии
5. **Состояния**: Используйте FSM для сложных диалогов
6. **Middleware**: Выносите общую логику в middleware
7. **Конфигурация**: Настраивайте под нагрузку

## Следующий шаг

Изучите [stop_bot](../stop_bot/) для graceful shutdown и управления жизненным циклом бота.