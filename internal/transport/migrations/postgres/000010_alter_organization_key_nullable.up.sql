-- remove default value
ALTER TABLE organizations
    ALTER COLUMN secret_key DROP DEFAULT;