-- Создаем расширение для генерации UUID если понадобится
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем схему для хранения истории заданий
CREATE TABLE IF NOT EXISTS tasks_history (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    solution TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска по дате создания
CREATE INDEX IF NOT EXISTS idx_tasks_history_created_at ON tasks_history(created_at);