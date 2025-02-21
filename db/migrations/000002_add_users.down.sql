DROP INDEX IF EXISTS idx_accounts_owner_currency;
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS fk_accounts_owner;
DROP TABLE IF EXISTS users;