-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create folders table with range partitioning and sub-partitioning by user_id
CREATE TABLE folders (
    id BIGSERIAL NOT NULL,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_folder_id BIGINT,
    size BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, user_id)
) PARTITION BY RANGE (user_id);

-- Partition creation for folders(as an example for test task)
CREATE TABLE folders_part_1 PARTITION OF folders
    FOR VALUES FROM (0) TO (1000);

-- Create files table with range partitioning and sub-partitioning by user_id
CREATE TABLE files (
    id BIGSERIAL NOT NULL,
    folder_id BIGINT NOT NULL,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    s3_url TEXT NOT NULL,
    size BIGINT NOT NULL,
    transaction_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, user_id)
) PARTITION BY RANGE (user_id);

-- Partition creation for files
CREATE TABLE files_part_1 PARTITION OF files
    FOR VALUES FROM (0) TO (1000);

-- Create upload_transactions table
CREATE TABLE upload_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    folder_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraints
ALTER TABLE folders ADD CONSTRAINT fk_folders_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE upload_transactions ADD CONSTRAINT fk_upload_transactions_user
    FOREIGN KEY (user_id) REFERENCES users(id);


-- Indexes for efficient querying
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_folder_parent ON folders(parent_folder_id, user_id);
CREATE INDEX idx_file_folder ON files(folder_id, user_id);
CREATE INDEX idx_upload_transaction_status ON upload_transactions(status);

-- Insert dummy user
INSERT INTO users (username, email) VALUES ('dummy_user', 'dummy_user@example.com');

-- Insert parent folder for dummy user
INSERT INTO folders (user_id, name, parent_folder_id)
VALUES ((SELECT id FROM users WHERE username = 'dummy_user'), '/', NULL);
