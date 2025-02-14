ALTER TABLE tasks ADD COLUMN user_id INT;

-- Убедимся, что колонка добавлена перед добавлением внешнего ключа
ALTER TABLE tasks
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
