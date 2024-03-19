CREATE TABLE IF NOT EXISTS "organizations" (
  "id" uuid PRIMARY KEY  DEFAULT uuid_generate_v4(),
  "name" varchar(20),
  "owner_id" uuid,
  "status" organization_status,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "organization_employees" (
  "id" uuid PRIMARY KEY  DEFAULT uuid_generate_v4(),
  "org_id" uuid,
  "user_id" uuid NOT NULL,
  "updated_at" timestamp,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "deleted_by" uuid
);

ALTER TABLE "organizations" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "organization_employees" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "organization_employees" ADD FOREIGN KEY ("org_id") REFERENCES "organizations" ("id");
