ALTER TABLE "organizations"
  ALTER COLUMN "status" SET DEFAULT 'ACTIVE';


ALTER TABLE "organizations"
  ALTER COLUMN "status" SET NOT NULL;