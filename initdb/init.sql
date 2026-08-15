CREATE TABLE IF NOT EXISTS random_number (
    id INT PRIMARY KEY,
    number INT NOT NULL
);

INSERT INTO random_number (id, number)
VALUES (1, 42)
ON DUPLICATE KEY UPDATE id = id;

-- Config user to have metric access
GRANT PROCESS, REPLICATION CLIENT, SELECT
ON *.*
TO 'test'@'%';