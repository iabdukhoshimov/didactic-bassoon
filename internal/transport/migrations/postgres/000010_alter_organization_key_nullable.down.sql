-- set default value for key to 'not generated'
ALTER TABLE organizations ALTER COLUMN secret_key SET DEFAULT 'not generated';
