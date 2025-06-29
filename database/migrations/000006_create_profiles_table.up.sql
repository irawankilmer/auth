CREATE TABLE profiles (
    id VARCHAR(26) NOT NULL UNIQUE PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL UNIQUE,
    full_name VARCHAR(125) NOT NULL,
    address TEXT,
    gender ENUM('laki-laki', 'perempuan'),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;