-- Enforce organization consistency and business identifiers after the initial schema.

ALTER TABLE persons
  ADD CONSTRAINT uq_person_id_org UNIQUE (id, org_id);

CREATE UNIQUE INDEX uq_person_org_email
  ON persons (org_id, lower(org_email))
  WHERE org_email IS NOT NULL;

ALTER TABLE employment_records
  ADD CONSTRAINT fk_employment_person_same_org
  FOREIGN KEY (person_id, org_id) REFERENCES persons (id, org_id);

CREATE UNIQUE INDEX uq_current_employee_id_org
  ON employment_records (org_id, employee_id)
  WHERE is_current = TRUE AND employee_id IS NOT NULL;

ALTER TABLE employment_records
  ADD CONSTRAINT fk_employment_manager_same_org
  FOREIGN KEY (reports_to, org_id) REFERENCES persons (id, org_id);

ALTER TABLE status_change_events
  ADD CONSTRAINT fk_event_person_same_org
  FOREIGN KEY (person_id, org_id) REFERENCES persons (id, org_id);

ALTER TABLE transfer_records
  ADD CONSTRAINT fk_transfer_person_same_org
  FOREIGN KEY (person_id, org_id) REFERENCES persons (id, org_id);

CREATE OR REPLACE FUNCTION reject_event_mutation() RETURNS TRIGGER AS $$
BEGIN
  RAISE EXCEPTION 'status_change_events is append-only';
END; $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_events_immutable ON status_change_events;
CREATE TRIGGER trg_events_immutable
  BEFORE UPDATE OR DELETE ON status_change_events
  FOR EACH ROW EXECUTE FUNCTION reject_event_mutation();
