CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Тестовые продукты
INSERT INTO products (name, description, price, stock, category) VALUES
('iPhone 15 Pro', 'Новейший смартфон от Apple', 99999.99, 10, 'Смартфоны'),
('MacBook Pro M3', 'Ноутбук с процессором M3', 149999.99, 5, 'Ноутбуки'),
('AirPods Pro', 'Беспроводные наушники с шумоподавлением', 24999.99, 20, 'Аксессуары'),
('Apple Watch Series 9', 'Умные часы с новейшими функциями', 39999.99, 15, 'Гаджеты'),
('iPad Pro', 'Планшет с дисплеем Liquid Retina', 79999.99, 8, 'Планшеты'); 