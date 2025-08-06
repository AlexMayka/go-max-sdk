# FSM Bot (Finite State Machine)

## Описание

Демонстрирует мощную систему состояний (FSM) на примере бота заказа пиццы. Показывает как создать многошаговый диалог с сохранением данных между этапами.

## Функциональность

- 🍕 Многошаговый процесс заказа пиццы
- 🔄 Управление состояниями пользователя
- 💾 Сохранение данных между шагами
- ✅ Валидация пользовательского ввода
- 🚪 Выход из состояния после завершения

## Схема состояний

```mermaid
graph TD
    A[Начало] --> B[/start - Приветствие]
    B --> C[/order - Выбор пиццы]
    C --> D[choosing_pizza]
    D --> E[choosing_size]
    E --> F[awaiting_address]
    F --> G[Заказ подтвержден]
    G --> A
```

## Состояния и переходы

### 1. Начальное состояние (нет состояния)
```go
bot.OnCommand("/start", startHandler)  // Приветствие
bot.OnCommand("/order", orderHandler)  // Начать заказ
```

### 2. Выбор пиццы (`choosing_pizza`)
```go
bot.UseState("choosing_pizza").Any(pizzaChoiceHandler)
```
- Пользователь вводит: `1`, `2`, или `3`
- Сохраняется: `pizza` → "Margherita"/"Pepperoni"/"Hawaiian"
- Переход в: `choosing_size`

### 3. Выбор размера (`choosing_size`)
```go
bot.UseState("choosing_size").Any(sizeChoiceHandler)
```
- Пользователь вводит: `1`, `2`, или `3`
- Сохраняется: `size` → "Small"/"Medium"/"Large"
- Сохраняется: `price` → "$10"/"$15"/"$20"
- Переход в: `awaiting_address`

### 4. Ввод адреса (`awaiting_address`)
```go
bot.UseState("awaiting_address").Any(addressHandler)
```
- Валидация: минимум 10 символов
- Формирование итогового заказа
- Очистка состояния: `ctx.SetState("")`

## Работа с данными

### Сохранение данных
```go
ctx.SetData("pizza", "Margherita")
ctx.SetData("size", "Large")
ctx.SetData("price", "$20")
```

### Получение данных
```go
pizza, _ := ctx.GetData("pizza")
size, _ := ctx.GetData("size")
price, _ := ctx.GetData("price")
```

## Как тестировать

1. Запустите бота: `go run main.go`
2. Начните диалог:

| Шаг | Команда/Сообщение | Результат |
|-----|-------------------|-----------|
| 1 | `/start` | Приветствие |
| 2 | `/order` | Меню выбора пиццы |
| 3 | `2` | Выбрана Pepperoni, меню размеров |
| 4 | `3` | Выбран Large, запрос адреса |
| 5 | `Moscow, Red Square 1` | Подтверждение заказа ✅ |

### Валидация ошибок
- Неправильный выбор: `"Please choose 1, 2, or 3"`
- Короткий адрес: `"Please enter a complete address"`

## Ключевые API

### Управление состояниями
```go
ctx.SetState("new_state")        // Установить состояние
currentState := ctx.GetState()   // Получить текущее состояние
ctx.SetState("")                 // Очистить состояние
```

### Работа с данными
```go
ctx.SetData("key", "value")      // Сохранить данные
value, exists := ctx.GetData("key")  // Получить данные
```

### Условная маршрутизация
```go
bot.UseState("state_name").Any(handler)      // Только в этом состоянии
bot.UseState("state_name").OnCommand("/cmd", handler)  // Команда + состояние
```

## Преимущества FSM

- ✅ **Структурированность**: Четкая логика переходов
- ✅ **Изоляция**: Каждое состояние независимо
- ✅ **Валидация**: Контроль правильности ввода
- ✅ **Персистентность**: Данные сохраняются между сообщениями
- ✅ **Масштабируемость**: Легко добавлять новые состояния

## Следующий шаг

Изучите [context_bot](../context_bot/) для работы с контекстом и расширенными данными.