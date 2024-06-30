BEGIN;
-- Insert initial data into users table
INSERT INTO users (passportNumber, surname, name, patronymic, address) VALUES
('1234567890', 'Ivanov', 'Ivan', 'Ivanovich', '123 Main St, Moscow'),
('2345678901', 'Petrov', 'Petr', 'Petrovich', '456 Elm St, St. Petersburg'),
('3456789012', 'Sidorov', 'Sidr', 'Sidorovich', '789 Pine St, Novosibirsk');
-- Insert initial data into tasks table with correlated subqueries to determine user_id
INSERT INTO tasks (userid, taskname, content, starttime, endtime) VALUES
((SELECT id FROM users WHERE passportNumber = '1234567890'), 'Task 1', 'Content 1', '2024-06-01 08:00:00', '2024-06-01 10:00:00'),
((SELECT id FROM users WHERE passportNumber = '1234567890'), 'Task 2', 'Content 2', '2024-06-02 09:00:00', '2024-06-02 12:00:00'),
((SELECT id FROM users WHERE passportNumber = '2345678901'), 'Task 3', 'Content 3', '2024-06-01 08:00:00', '2024-06-01 09:30:00'),
((SELECT id FROM users WHERE passportNumber = '2345678901'), 'Task 4', 'Content 4', '2024-06-03 10:00:00', '2024-06-03 13:00:00'),
((SELECT id FROM users WHERE passportNumber = '3456789012'), 'Task 5', 'Content 5', '2024-06-01 07:00:00', '2024-06-01 08:00:00'),
((SELECT id FROM users WHERE passportNumber = '3456789012'), 'Task 6', 'Content 6', '2024-06-04 11:00:00', '2024-06-04 12:00:00');
COMMIT;