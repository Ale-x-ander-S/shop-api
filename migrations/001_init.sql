CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    stock INTEGER DEFAULT 0,
    category VARCHAR(255),
    image_url VARCHAR(255)
);

-- Тестовые данные
INSERT INTO products (name, description, price, stock, category, image_url) 
VALUES 
    ('Test Product 1', 'Description 1', 99.99, 10, 'Electronics', 'https://example.com/img1.jpg'),
    ('Test Product 2', 'Description 2', 149.99, 5, 'Books', 'https://example.com/img2.jpg'),
    ('Test Product 3', 'Description 3', 199.99, 15, 'Clothing', 'https://example.com/img3.jpg')
ON CONFLICT DO NOTHING;

-- Тестовые продукты
INSERT INTO products (name, description, price, stock, category) VALUES
('iPhone 15 Pro', 'Новейший смартфон от Apple', 99999.99, 10, 'Смартфоны'),
('MacBook Pro M3', 'Ноутбук с процессором M3', 149999.99, 5, 'Ноутбуки'),
('AirPods Pro', 'Беспроводные наушники с шумоподавлением', 24999.99, 20, 'Аксессуары'),
('Apple Watch Series 9', 'Умные часы с новейшими функциями', 39999.99, 15, 'Гаджеты'),
('iPad Pro', 'Планшет с дисплеем Liquid Retina', 79999.99, 8, 'Планшеты'); 