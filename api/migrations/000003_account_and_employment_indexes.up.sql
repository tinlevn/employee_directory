-- Enforce 1:1 account per person
CREATE UNIQUE INDEX idx_unique_account_per_person ON person_accounts(person_id);

-- Partial index for high-traffic current employment lookup
CREATE INDEX idx_employment_current_person ON employment_records(person_id) WHERE is_current = TRUE;
