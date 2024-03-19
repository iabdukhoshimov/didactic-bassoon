
CREATE TABLE IF NOT EXISTS "refresh_token"(
   "id" uuid PRIMARY KEY  DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "users" (id),
    refresh_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
)
