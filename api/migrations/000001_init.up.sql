-- 000001_init.up.sql
-- Employee Directory — trimmed schema (removed PII: ethnicity/religion/blood_type, maiden_name, pronouns, nationality, phone_secondary, linkedin, website, inline emergency_contact -> separate table; removed cost_center/division/business_unit)
-- Postgres 16+, UUID gen_random_uuid() requires pgcrypto

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- organizations
CREATE TABLE organizations (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name           VARCHAR(255) NOT NULL,
  type           VARCHAR(50)  NOT NULL CHECK (type IN ('company', 'ngo', 'government', 'other')),
  country        VARCHAR(100),
  timezone       VARCHAR(100),
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- persons (trimmed)
CREATE TABLE persons (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id                UUID NOT NULL REFERENCES organizations(id),

  first_name            VARCHAR(100) NOT NULL,
  middle_name           VARCHAR(100),
  last_name             VARCHAR(100) NOT NULL,
  preferred_name        VARCHAR(100),
  date_of_birth         DATE,
  gender                VARCHAR(50) CHECK (gender IN ('male','female','non-binary','prefer-not-to-say')),
  profile_photo_url     TEXT,

  personal_email        VARCHAR(255),
  org_email             VARCHAR(255),
  phone_primary         VARCHAR(50),

  address_line_1        VARCHAR(255),
  address_line_2        VARCHAR(255),
  city                  VARCHAR(100),
  state_province        VARCHAR(100),
  postal_code           VARCHAR(20),
  country               VARCHAR(100),
  is_international      BOOLEAN DEFAULT FALSE,

  is_active             BOOLEAN NOT NULL DEFAULT TRUE,
  archived_at           TIMESTAMPTZ,
  archive_reason        VARCHAR(100),
  source                VARCHAR(100),
  notes                 TEXT,
  tags                  TEXT[],

  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by            UUID,
  updated_by            UUID
);

CREATE INDEX idx_persons_org ON persons(org_id);
CREATE INDEX idx_persons_active ON persons(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_persons_name ON persons(last_name, first_name);
CREATE INDEX idx_persons_email ON persons(org_email);
CREATE INDEX idx_persons_tags ON persons USING GIN (tags);

-- emergency_contacts — separate entity, optional FK to persons, 1 per person (extensible)
CREATE TABLE emergency_contacts (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  name            VARCHAR(200) NOT NULL,
  phone           VARCHAR(50),
  email           VARCHAR(255),
  relationship    VARCHAR(100),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (person_id)
);
CREATE INDEX idx_emergency_person ON emergency_contacts(person_id);

-- employment_records SCD2 (trimmed)
CREATE TABLE employment_records (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id           UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  org_id              UUID NOT NULL REFERENCES organizations(id),

  employee_id         VARCHAR(100),
  job_title           VARCHAR(200),
  job_level           VARCHAR(100),
  employment_status   VARCHAR(50) CHECK (employment_status IN ('full-time', 'part-time', 'contract', 'intern','freelance', 'probationary', 'on-leave')),
  employment_type     VARCHAR(50) CHECK (employment_type IN ('permanent', 'fixed-term', 'casual')),
  work_arrangement    VARCHAR(50) CHECK (work_arrangement IN ('on-site', 'remote', 'hybrid')),

  department          VARCHAR(200),
  team                VARCHAR(200),
  office_location     VARCHAR(200),
  desk_number         VARCHAR(50),
  reports_to          UUID REFERENCES persons(id),

  salary_amount       BIGINT,
  salary_currency     VARCHAR(10),
  pay_frequency       VARCHAR(20) CHECK (pay_frequency IN ('monthly', 'bi-weekly', 'weekly')),
  hourly_rate         NUMERIC(10,2),

  hire_date           DATE,
  probation_end_date  DATE,
  contract_start_date DATE,
  contract_end_date   DATE,

  valid_from          DATE NOT NULL,
  valid_to            DATE,
  is_current          BOOLEAN NOT NULL DEFAULT TRUE,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by          UUID,
  updated_by          UUID,

  CONSTRAINT chk_valid_dates CHECK (valid_to IS NULL OR valid_to > valid_from)
);

CREATE INDEX idx_employment_person ON employment_records(person_id);
CREATE UNIQUE INDEX idx_employment_current ON employment_records(person_id) WHERE is_current = TRUE;
CREATE INDEX idx_employment_dates ON employment_records(valid_from, valid_to);
CREATE INDEX idx_employment_org_dept ON employment_records(org_id, department);

-- status_change_events append-only
CREATE TABLE status_change_events (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  org_id          UUID NOT NULL REFERENCES organizations(id),

  event_type      VARCHAR(50) NOT NULL CHECK (event_type IN ('HIRED', 'REHIRED', 'PROMOTED', 'DEMOTED','TRANSFERRED', 'SECONDMENT_START', 'SECONDMENT_END','SALARY_CHANGE', 'TITLE_CHANGE','RESIGNED', 'TERMINATED', 'LAID_OFF', 'RETIRED','ON_LEAVE_START', 'ON_LEAVE_END','CONTRACT_RENEWED', 'CONTRACT_EXPIRED','ACTIVATED', 'DEACTIVATED', 'RECORD_UPDATED')),
  context         VARCHAR(20) NOT NULL CHECK (context IN ('employment', 'general')) DEFAULT 'employment',

  from_status     VARCHAR(100),
  to_status       VARCHAR(100),
  from_department VARCHAR(200),
  to_department   VARCHAR(200),
  from_title      VARCHAR(200),
  to_title        VARCHAR(200),
  from_location   VARCHAR(200),
  to_location     VARCHAR(200),

  reason          TEXT,
  reason_code     VARCHAR(100),
  is_voluntary    BOOLEAN,

  effective_date  DATE NOT NULL,
  recorded_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

  initiated_by    UUID REFERENCES persons(id),
  approved_by     UUID REFERENCES persons(id),
  witnessed_by    UUID REFERENCES persons(id),

  linked_record_id   UUID,
  linked_record_type VARCHAR(50),

  document_urls   TEXT[],
  notes           TEXT,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by      UUID
);

CREATE INDEX idx_events_person ON status_change_events(person_id);
CREATE INDEX idx_events_type ON status_change_events(event_type);
CREATE INDEX idx_events_date ON status_change_events(effective_date);
CREATE INDEX idx_events_org ON status_change_events(org_id);
CREATE INDEX idx_events_recorded ON status_change_events(recorded_at DESC);

CREATE OR REPLACE FUNCTION reject_event_mutation() RETURNS TRIGGER AS $$
BEGIN
  RAISE EXCEPTION 'status_change_events is append-only';
END; $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_events_immutable
  BEFORE UPDATE OR DELETE ON status_change_events
  FOR EACH ROW EXECUTE FUNCTION reject_event_mutation();

-- transfer_records
CREATE TABLE transfer_records (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id        UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  org_id           UUID NOT NULL REFERENCES organizations(id),

  transfer_type    VARCHAR(50) CHECK (transfer_type IN ('INTERNAL', 'INTER-CAMPUS', 'SECONDMENT','PROMOTION', 'DEMOTION', 'LATERAL')),
  from_department  VARCHAR(200),
  to_department    VARCHAR(200),
  from_location    VARCHAR(200),
  to_location      VARCHAR(200),
  from_manager_id  UUID REFERENCES persons(id),
  to_manager_id    UUID REFERENCES persons(id),
  from_title       VARCHAR(200),
  to_title         VARCHAR(200),
  effective_date   DATE NOT NULL,
  reason           TEXT,
  approved_by      UUID REFERENCES persons(id),
  notes            TEXT,

  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by       UUID
);
CREATE INDEX idx_transfers_person ON transfer_records(person_id);
CREATE INDEX idx_transfers_date ON transfer_records(effective_date);

-- headcount_snapshots
CREATE TABLE headcount_snapshots (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id           UUID NOT NULL REFERENCES organizations(id),
  snapshot_date    DATE NOT NULL,
  department       VARCHAR(200),
  location         VARCHAR(200),
  active_count     INT NOT NULL DEFAULT 0,
  new_entries      INT NOT NULL DEFAULT 0,
  exits            INT NOT NULL DEFAULT 0,
  net_change       INT GENERATED ALWAYS AS (new_entries - exits) STORED,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE NULLS NOT DISTINCT (org_id, snapshot_date, department, location)
);
CREATE INDEX idx_snapshots_date ON headcount_snapshots(snapshot_date);
CREATE INDEX idx_snapshots_org ON headcount_snapshots(org_id);

-- supporting tables
CREATE TABLE person_skills (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  skill_name      VARCHAR(200) NOT NULL,
  proficiency     VARCHAR(50) CHECK (proficiency IN ('beginner','intermediate','advanced','expert')),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_skills_person ON person_skills(person_id);

CREATE TABLE person_languages (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  language        VARCHAR(100) NOT NULL,
  proficiency     VARCHAR(10) CHECK (proficiency IN ('A1','A2','B1','B2','C1','C2','native')),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE person_documents (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  document_type   VARCHAR(100),
  file_url        TEXT NOT NULL,
  expiry_date     DATE,
  uploaded_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  uploaded_by     UUID
);

CREATE TABLE person_certifications (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  cert_name       VARCHAR(200) NOT NULL,
  issuing_body    VARCHAR(200),
  issued_date     DATE,
  expiry_date     DATE,
  credential_url  TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE employee_benefits (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id       UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  benefit_type    VARCHAR(100),
  provider        VARCHAR(200),
  start_date      DATE,
  end_date        DATE,
  details         JSONB,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE person_accounts (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  person_id          UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
  username           VARCHAR(200) UNIQUE NOT NULL,
  password_hash      TEXT NOT NULL,
  role               VARCHAR(50) CHECK (role IN ('admin','manager','staff','read-only')),
  permissions        TEXT[],
  is_active          BOOLEAN NOT NULL DEFAULT TRUE,
  last_login         TIMESTAMPTZ,
  two_factor_enabled BOOLEAN DEFAULT FALSE,
  account_locked     BOOLEAN DEFAULT FALSE,
  sso_provider       VARCHAR(50),
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_accounts_person ON person_accounts(person_id);
CREATE INDEX idx_accounts_username ON person_accounts(username);

-- updated_at trigger
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_org_updated BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_person_updated BEFORE UPDATE ON persons FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_emergency_updated BEFORE UPDATE ON emergency_contacts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_employment_updated BEFORE UPDATE ON employment_records FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_account_updated BEFORE UPDATE ON person_accounts FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- views
CREATE VIEW active_employees AS
SELECT p.id, p.org_id, p.first_name, p.last_name, p.org_email, p.is_active,
       e.job_title, e.job_level, e.department, e.team, e.office_location, e.employment_status,
       e.hire_date, e.reports_to, e.salary_amount, e.salary_currency
FROM persons p
JOIN employment_records e ON e.person_id = p.id AND e.is_current = TRUE
WHERE p.is_active = TRUE;

CREATE VIEW person_event_timeline AS
SELECT e.person_id, e.event_type, e.context, e.from_status, e.to_status,
       e.from_department, e.to_department, e.from_title, e.to_title,
       e.reason, e.is_voluntary, e.effective_date, e.recorded_at,
       p_init.first_name || ' ' || p_init.last_name AS initiated_by_name,
       p_appr.first_name || ' ' || p_appr.last_name AS approved_by_name
FROM status_change_events e
LEFT JOIN persons p_init ON p_init.id = e.initiated_by
LEFT JOIN persons p_appr ON p_appr.id = e.approved_by
ORDER BY e.effective_date DESC, e.recorded_at DESC;

-- seed default org for local dev
INSERT INTO organizations (id, name, type, country, timezone) VALUES
('00000000-0000-0000-0000-000000000001', 'Default Org', 'company', 'USA', 'America/New_York')
ON CONFLICT (id) DO NOTHING;
