-- Drop the files table
DROP TABLE IF EXISTS files;

-- Drop the folders table
DROP TABLE IF EXISTS folders;

-- Drop the upload_transactions table
DROP TABLE IF EXISTS upload_transactions;

-- Drop the users table
DROP TABLE IF EXISTS users;

-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_folder_user;
DROP INDEX IF EXISTS idx_file_folder;
DROP INDEX IF EXISTS idx_file_user;
DROP INDEX IF EXISTS idx_upload_transaction_status;
