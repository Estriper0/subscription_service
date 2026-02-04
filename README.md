# Subscription service

Сервис для управления онлайн-подписками пользователей. 

---

# Шаги по запуску

1. **Клонируй репозиторий и перейдите в папку**:
   ```
   git clone https://github.com/Estriper0/subscription_service.git
   cd subscription_service
   ```
2. Создайте файл `.env` и настройте переменные окружения:
   ```env
    ENV=local

    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=12345
    DB_NAME=postgres
   ```

3. **Запусти с помощью Make**:
   ```
   make up
   ```

---

По пути <localhost:8080/swagger/index.html> можно ознакомиться с документацией

