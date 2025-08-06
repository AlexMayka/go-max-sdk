# BotContext

Полный справочник по работе с контекстом бота.

## Структура BotContext

```go
type BotContext struct {
    Ctx    context.Context  // Go context для отмены операций
    Cancel context.CancelFunc // Функция отмены
    
    // Информация о пользователе
    UserID    int64   // ID пользователя
    ChatID    int64   // ID чата (0 для личных сообщений)
    FirstName string  // Имя пользователя
    LastName  string  // Фамилия пользователя  
    Username  string  // Username пользователя
    
    // Данные обновления
    Update *models.Update  // Полное обновление от API
    Text   string         // Текст сообщения
    
    // Внутренние компоненты
    Client APIClient  // HTTP клиент для API запросов
    FSM    FSM       // Конечный автомат для состояний
}
```

## Методы отправки сообщений

### Reply()
Отправляет текстовое сообщение.

```go
func (ctx *BotContext) Reply(text string)
```

**Примеры:**
```go
ctx.Reply("Привет!")
ctx.Reply("Ваш баланс: 100₽")
ctx.Reply("❌ Ошибка: неверная команда")
```

### ReplyWithKeyboard()
Отправляет сообщение с inline клавиатурой.

```go
func (ctx *BotContext) ReplyWithKeyboard(text string, keyboard *models.InlineKeyboard) error
```

**Пример:**
```go
yesPayload := "yes"
noPayload := "no"

keyboard := &models.InlineKeyboard{
    Buttons: [][]models.InlineKeyboardButton{
        {
            {
                Type: models.InlineKeyboardButtonTypeCallback,
                Text: "✅ Да", 
                Payload: &yesPayload,
            },
            {
                Type: models.InlineKeyboardButtonTypeCallback,
                Text: "❌ Нет", 
                Payload: &noPayload,
            },
        },
    },
}
ctx.ReplyWithKeyboard("Подтвердить действие?", keyboard)
```

## Методы управления сообщениями

### EditMessage()
Редактирует существующее сообщение.

```go
func (ctx *BotContext) EditMessage(messageID, newText string) error
```

**Пример:**
```go
err := ctx.EditMessage("msg_12345", "Обновленный текст")
if err != nil {
    ctx.Reply("Ошибка редактирования")
}
```

### DeleteMessage()
Удаляет сообщение по ID.

```go
func (ctx *BotContext) DeleteMessage(messageID string) error
```

**Пример:**
```go
err := ctx.DeleteMessage("msg_12345")
if err != nil {
    log.Printf("Failed to delete message: %v", err)
}
```

## Методы работы с состояниями (FSM)

### SetState()
Устанавливает состояние для пользователя.

```go
func (ctx *BotContext) SetState(state string) bool
```

**Примеры:**
```go
ctx.SetState("waiting_name")     // Ожидание имени
ctx.SetState("choosing_size")    // Выбор размера
ctx.SetState("")                 // Сброс состояния
```

### GetState()
Получает текущее состояние пользователя.

```go
func (ctx *BotContext) GetState() (string, bool)
```

**Пример:**
```go
state, exists := ctx.GetState()
if exists {
    ctx.Reply("Ваше состояние: " + state)
} else {
    ctx.Reply("Состояние не установлено")
}
```

## Методы работы с данными (FSM)

### SetData()
Сохраняет данные для пользователя.

```go
func (ctx *BotContext) SetData(key, value string) bool
```

**Примеры:**
```go
ctx.SetData("name", "Иван")
ctx.SetData("email", "ivan@example.com")
ctx.SetData("cart_total", "1500")
```

### GetData()
Получает сохраненные данные пользователя.

```go
func (ctx *BotContext) GetData(key string) (string, bool)
```

**Пример:**
```go
name, exists := ctx.GetData("name")
if exists {
    ctx.Reply("Привет, " + name + "!")
} else {
    ctx.Reply("Имя не указано")
}
```

## Доступ к информации о пользователе

### Базовая информация

```go
// ID пользователя (всегда есть)
userID := ctx.UserID

// Полное имя
fullName := ctx.FirstName
if ctx.LastName != "" {
    fullName += " " + ctx.LastName
}

// Username (может быть пустым)
if ctx.Username != "" {
    ctx.Reply("@" + ctx.Username + ", привет!")
}
```

### Определение типа чата

```go
if ctx.ChatID != 0 {
    ctx.Reply("Это групповой чат")
} else {
    ctx.Reply("Это личное сообщение")
}
```

## Работа с Update

### Доступ к полному обновлению

```go
update := ctx.Update

// Информация о сообщении
if update.Message != nil {
    messageID := update.Message.Body.Mid
    timestamp := update.Message.Timestamp
}

// Callback data (для inline кнопок)
if update.UpdateType == models.UpdateTypeMessageCallback {
    data := update.Message.Body.Text
    ctx.Reply("Вы нажали: " + data)
}
```

### Типы обновлений

```go
switch {
case ctx.Update.Message != nil:
    ctx.Reply("Получено сообщение: " + ctx.Text)
    
case ctx.Update.UpdateType == models.UpdateTypeMessageCallback:
    ctx.Reply("Нажата кнопка: " + ctx.Text)
    
case ctx.Update.UpdateType == models.UpdateTypeMessageEdited:
    ctx.Reply("Сообщение отредактировано")
}
```

## Практические примеры

### Регистрация пользователя

```go
func registerStart(ctx *BotContext) {
    ctx.Reply("Регистрация. Введите ваше имя:")
    ctx.SetState("reg_name")
}

func registerName(ctx *BotContext) {
    name := ctx.Text
    if len(name) < 2 {
        ctx.Reply("Имя слишком короткое. Повторите:")
        return
    }
    
    ctx.SetData("name", name)
    ctx.Reply("Имя: " + name + "\nВведите email:")
    ctx.SetState("reg_email")
}

func registerEmail(ctx *BotContext) {
    email := ctx.Text
    name, _ := ctx.GetData("name")
    
    ctx.Reply("Регистрация завершена!\n👤 " + name + "\n📧 " + email)
    ctx.SetData("registered", "true")
    ctx.SetState("") // Сброс состояния
}
```

### Корзина покупок

```go
func addToCart(ctx *BotContext) {
    item := strings.TrimPrefix(ctx.Text, "/add ")
    
    // Получаем текущую корзину
    cart, exists := ctx.GetData("cart")
    if !exists {
        cart = ""
    }
    
    // Добавляем товар
    if cart != "" {
        cart += "," + item
    } else {
        cart = item
    }
    
    ctx.SetData("cart", cart)
    ctx.Reply("✅ Добавлено: " + item)
}

func showCart(ctx *BotContext) {
    cart, exists := ctx.GetData("cart")
    if !exists || cart == "" {
        ctx.Reply("🛒 Корзина пуста")
        return
    }
    
    items := strings.Split(cart, ",")
    message := "🛒 Ваша корзина:\n"
    for i, item := range items {
        message += fmt.Sprintf("%d. %s\n", i+1, item)
    }
    
    ctx.Reply(message)
}
```

### Авторизация админа

```go
func adminLogin(ctx *BotContext) {
    password := strings.TrimPrefix(ctx.Text, "/admin ")
    
    if password == "secret123" {
        ctx.SetData("admin", "true")
        ctx.Reply("✅ Админ доступ получен")
    } else {
        ctx.Reply("❌ Неверный пароль")
    }
}

func adminOnly(ctx *BotContext) {
    admin, exists := ctx.GetData("admin")
    if !exists || admin != "true" {
        ctx.Reply("🔒 Нужны права админа")
        return
    }
    
    ctx.Reply("👑 Добро пожаловать, админ!")
}
```

## Методы Bot (вне контекста)

Помимо методов BotContext, есть методы Bot для отправки сообщений вне обработчиков.

### Start()
Запускает бота. Блокирующий вызов.

```go
func (b *Bot) Start() error
```

**Пример:**
```go
bot := maxsdk.NewBot("token")
bot.OnCommand("/start", startHandler)

if err := bot.Start(); err != nil {
    log.Fatal(err)
}
```

### Stop()
Graceful остановка бота.

```go
func (b *Bot) Stop() error
```

**Пример:**
```go
err := bot.Stop()
if err != nil {
    log.Printf("Error stopping bot: %v", err)
}
```

### SendText()
Отправляет текстовое сообщение пользователю или в чат.

```go
func (b *Bot) SendText(userID, chatID int64, text string) error
```

**Примеры:**
```go
// В личные сообщения
err := bot.SendText(12345, 0, "Привет!")

// В групповой чат  
err := bot.SendText(0, 67890, "Сообщение в группу")
```

### SendTextToUser()
Отправляет текстовое сообщение пользователю.

```go
func (b *Bot) SendTextToUser(userID int64, text string) error
```

**Пример:**
```go
err := bot.SendTextToUser(12345, "Личное сообщение")
```

### SendTextToChat()
Отправляет текстовое сообщение в чат.

```go
func (b *Bot) SendTextToChat(chatID int64, text string) error
```

**Пример:**
```go
err := bot.SendTextToChat(67890, "Сообщение в чат")
```

### SendTextWithKeyboard()
Отправляет сообщение с inline клавиатурой.

```go
func (b *Bot) SendTextWithKeyboard(userID, chatID int64, text string, keyboard *models.InlineKeyboard) error
```

**Пример:**
```go
yesPayload := "yes"
noPayload := "no"

keyboard := &models.InlineKeyboard{
    Buttons: [][]models.InlineKeyboardButton{
        {{
            Type: models.InlineKeyboardButtonTypeCallback,
            Text: "Да",
            Payload: &yesPayload,
        }},
        {{
            Type: models.InlineKeyboardButtonTypeCallback,
            Text: "Нет",
            Payload: &noPayload,
        }},
    },
}

err := bot.SendTextWithKeyboard(12345, 0, "Вопрос?", keyboard)
```

### EditMessage() (Bot метод)
Редактирует сообщение по ID.

```go
func (b *Bot) EditMessage(messageID, newText string) error
```

**Пример:**
```go
err := bot.EditMessage("msg_123", "Обновленный текст")
```

### DeleteMessage() (Bot метод)  
Удаляет сообщение по ID.

```go
func (b *Bot) DeleteMessage(messageID string) error
```

**Пример:**
```go
err := bot.DeleteMessage("msg_123")
```

## Best Practices

### 1. Проверка данных
```go
// ❌ Плохо
name, _ := ctx.GetData("name")
ctx.Reply("Привет, " + name)

// ✅ Хорошо
name, exists := ctx.GetData("name")
if !exists {
    ctx.Reply("Сначала зарегистрируйтесь")
    return
}
ctx.Reply("Привет, " + name)
```

### 2. Валидация ввода
```go
func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func handleEmail(ctx *BotContext) {
    email := ctx.Text
    if !validateEmail(email) {
        ctx.Reply("❌ Неверный формат email")
        return
    }
    
    ctx.SetData("email", email)
    ctx.Reply("✅ Email сохранен")
}
```

### 3. Очистка состояний
```go
func completeProcess(ctx *BotContext) {
    // Завершаем процесс
    ctx.Reply("✅ Готово!")
    
    // Очищаем временные данные
    ctx.SetData("temp_data", "")
    ctx.SetState("") // Сброс состояния
}
```

### 4. Обработка ошибок
```go
func sendMessage(ctx *BotContext) {
    err := ctx.EditMessage("msg_123", "Новый текст")
    if err != nil {
        ctx.Reply("⚠️ Не удалось отредактировать сообщение")
        return
    }
    
    ctx.Reply("✅ Сообщение обновлено")
}
```