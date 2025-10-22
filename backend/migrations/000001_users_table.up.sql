CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE CHECK(POSITION('@' IN email) > 1),
    username VARCHAR(64) NOT NULL UNIQUE CHECK(LENGTH(username) >= 5 AND LENGTH(username) < 64),
    hashed_password VARCHAR(72) NOT NULL, 
    bio TEXT, 
    profile_status VARCHAR(50),
    is_admin BOOLEAN DEFAULT FALSE, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);
