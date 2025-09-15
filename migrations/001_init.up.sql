-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Создание таблицы страниц
CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT now()
);

-- Индекс для быстрого поиска страниц по пользователю
CREATE INDEX IF NOT EXISTS idx_pages_user_id ON pages(user_id);