CREATE TYPE "roles" AS ENUM (
  'OWNER',
  'ADMIN',
  'USER',
  'EMPLOYEE'
);

CREATE TYPE "user_status" AS ENUM (
  'ACTIVE',
  'PENDING',
  'FREE_TRIAL',
  'INACTIVE',
  'DELETED',
  'BLOCKED'
);

CREATE TYPE "organization_status" AS ENUM (
  'ACTIVE',
  'PENDING',
  'FREE_TRIAL',
  'INACTIVE',
  'DELETED',
  'BLOCKED'
);
