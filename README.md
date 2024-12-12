
## Запуск

1. Создайте .env файл в корневом каталоге и добавьте следующие значения:

    ```
    LOG_LEVEL=
    
    ACCESS_TOKEN_TTL=
    REFRESH_TOKEN_TTL=
    JWT_SIGNING_KEY=
    
    SMTP_HOST=
    SMTP_PORT=
    FROM_EMAIL=
    FROM_PASSWORD=
    ```

2. Создайте образ приложения

    ```
   docker build -t auth-service .
   ```
3. Запустите проект

    ```
   doker-compose up
   ```
4. Сервис работает с пользователями, которые есть в базе данных. Добавление пользователей
происходит вместе с созданием базы данных в скрипте assets/init.sql. 
Для корректной отправки warning сообщений необходимо изменить почту.