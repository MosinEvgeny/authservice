CREATE TABLE refresh_tokens (
    user_id UUID PRIMARY KEY,
    hashed_refresh_token TEXT NOT NULL,
    access_token TEXT NOT NULL,
	client_ip TEXT NOT NULL
);

CREATE TABLE users (
	id UUID PRIMARY KEY,
	email TEXT NOT NULL UNIQUE
);