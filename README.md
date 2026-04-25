# Wallet App

Приложение для управления балансом кошельков.

## Запросы

### POST /api/v1/wallet

Создает кошелек или обновляет баланс.

**Тело запроса:**
```json
{
  "valletId": "UUID",
  "operationType": "DEPOSIT or WITHDRAW",
  "amount": "int"
}
```
**Пример запроса:**
```json
{
  "valletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

### GET /api/v1/wallets/{WALLET_UUID}

Возвращает баланс кошелька или ошибку, если кошелька не существует.


## Запуск

Для запуска приложения используйте
```bash
docker-compose up
```

Приложение запускает два Docker контейнера
- Go 
- PostgreSQL

Переменные среды хранятся в `config.env`.
