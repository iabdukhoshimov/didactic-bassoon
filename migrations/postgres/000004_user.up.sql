CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY  DEFAULT uuid_generate_v4(),
  "first_name" varchar(30) NOT NULL,
  "last_name" varchar(30) NOT NULL,
  "email" varchar(150) UNIQUE NOT NULL,
  "password" varchar(50) NOT NULL,
  "role" roles NOT NULL DEFAULT 'USER',
  "user_status" user_status,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

