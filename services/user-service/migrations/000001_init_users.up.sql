CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);


-- Migration commands:
-- To create the migration, run:
--   migrate create -ext sql -dir services/user-service/migrations -seq init_users
-- To apply the migration, run:
--   migrate -path services/user-service/migrations -database "your_database_url" up
-- If already migrated to some version, use:
--   migrate -path services/user-service/migrations -database "your_database_url" force <version_number>
-- To rollback the migration, run:
--   migrate -path services/user-service/migrations -database "your_database_url" down