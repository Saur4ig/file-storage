-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create folders table
CREATE TABLE folders (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    parent_folder_id BIGINT DEFAULT NULL,
    size BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, parent_folder_id, name),
    FOREIGN KEY (parent_folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

-- Create upload_transactions table
CREATE TABLE upload_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    folder_id BIGINT NOT NULL REFERENCES folders(id),
    total_size BIGINT DEFAULT 0,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create files table
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    folder_id BIGINT NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    s3_url TEXT NOT NULL,
    size BIGINT NOT NULL,
    transaction_id BIGINT REFERENCES upload_transactions(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(folder_id, name)
);


-- Indexes for efficient querying
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_folder_user ON folders(user_id);
CREATE INDEX idx_file_folder ON files(folder_id);
CREATE INDEX idx_file_user ON files(user_id);
CREATE INDEX idx_upload_transaction_status ON upload_transactions(status);

-- Insert dummy user
INSERT INTO users (username, email) VALUES ('dummy_user', 'dummy_user@example.com');

-- Insert parent folder for dummy user
INSERT INTO folders (user_id, name) VALUES ((SELECT id FROM users WHERE username = 'dummy_user'), '/');