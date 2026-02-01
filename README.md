🤖 Telegram-каталог товаров
Простой, но функциональный Telegram-бот для просмотра каталога товаров с поддержкой изображений, навигации и связи с менеджером.
Идеально подходит для малого бизнеса, маркетплейсов или учебных проектов.

    💡 Бот написан на Go (Golang) и работает с любой SQL-СУБД, совместимой с интерфейсом database/sql. По умолчанию используется PostgreSQL, но легко адаптируется под MySQL, SQLite и другие.

🛠️ Что нужно для запуска
1. Создать Telegram-бота

    Напишите @BotFather
     → /newbot  
    Получите токен и сохраните его

2. Настроить СУБД

        -- Таблица товаров
       CREATE TABLE products (
          id SERIAL PRIMARY KEY,
          name VARCHAR(255) NOT NULL,
          price INTEGER NOT NULL
        );

        -- Таблица изображений (один товар — много изображений)
        CREATE TABLE product_images (
          id SERIAL PRIMARY KEY,
          product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
          image_url TEXT NOT NULL
        );
🔗 Все image_url должны быть публично доступны по HTTP/HTTPS!

3. Заполнить данные
   
Добавьте несколько товаров и изображений вручную или через скрипт.
4. Указать параметры подключения

В файле main.go замените все заглушки YOUR_* на реальные значения:

    const (
      host     = "localhost"        // YOUR_HOST
      port     = 5432               // YOUR_PORT
      user     = "postgres"         // YOUR_USER
      password = "your_password"    // YOUR_PASSWORD
      dbname   = "shop_db"          // YOUR_DB
    )

И не забудьте вставить токен бота:

    bot, err := tgbotapi.NewBotAPI("YOUR_TOKEN")

🧩 Совместимость с другими СУБД
Хотя проект изначально настроен под PostgreSQL, вы можете легко переключиться на другую СУБД:

    MySQL: добавьте драйвер github.com/go-sql-driver/mysql, используйте sql.Open("mysql", "...") и замените $1 → ? в запросах  
    SQLite: добавьте github.com/mattn/go-sqlite3, используйте sql.Open("sqlite3", "file.db")

📦 Структура проекта

    telegram-catalog-bot/
    ├── main.go              # Основной код бота
    ├── go.mod               # Зависимости
    ├── go.sum
    └── README.md            # 📄 Этот файл!

🎯 Возможности

    📋 Просмотр списка товаров из БД  
    🖼 Отображение изображений по публичным URL  
    💬 Кнопка «Связаться с менеджером» (ваша ссылка)  
    🔁 Навигация: выбрать другой товар / завершить  
    🧱 Простая архитектура — легко расширять

