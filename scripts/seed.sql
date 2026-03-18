-- Seed product data
INSERT INTO products (id, name, price, created_at) VALUES 
(gen_random_uuid(), 'iPhone 15 Pro', 999.99, NOW() - INTERVAL '10 days'),
(gen_random_uuid(), 'MacBook Air M3', 1099.00, NOW() - INTERVAL '9 days'),
(gen_random_uuid(), 'AirPods Pro 2', 249.00, NOW() - INTERVAL '8 days'),
(gen_random_uuid(), 'iPad Pro', 799.00, NOW() - INTERVAL '7 days'),
(gen_random_uuid(), 'Apple Watch Ultra', 799.00, NOW() - INTERVAL '6 days'),
(gen_random_uuid(), 'Sony WH-1000XM5', 348.00, NOW() - INTERVAL '5 days'),
(gen_random_uuid(), 'Samsung Galaxy S24', 799.99, NOW() - INTERVAL '4 days'),
(gen_random_uuid(), 'Google Pixel 8', 699.00, NOW() - INTERVAL '3 days'),
(gen_random_uuid(), 'Nintendo Switch', 299.00, NOW() - INTERVAL '2 days'),
(gen_random_uuid(), 'PS5 Slim', 449.00, NOW() - INTERVAL '1 day'),
(gen_random_uuid(), 'Xbox Series X', 499.00, NOW()),
(gen_random_uuid(), 'GoPro Hero 12', 399.00, NOW() + INTERVAL '1 hour');
