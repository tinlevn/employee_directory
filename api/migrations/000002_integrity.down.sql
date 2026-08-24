DROP TRIGGER IF EXISTS trg_events_immutable ON status_change_events;
DROP FUNCTION IF EXISTS reject_event_mutation();
ALTER TABLE transfer_records DROP CONSTRAINT IF EXISTS fk_transfer_person_same_org;
ALTER TABLE status_change_events DROP CONSTRAINT IF EXISTS fk_event_person_same_org;
ALTER TABLE employment_records DROP CONSTRAINT IF EXISTS fk_employment_manager_same_org;
DROP INDEX IF EXISTS uq_current_employee_id_org;
ALTER TABLE employment_records DROP CONSTRAINT IF EXISTS fk_employment_person_same_org;
DROP INDEX IF EXISTS uq_person_org_email;
ALTER TABLE persons DROP CONSTRAINT IF EXISTS uq_person_id_org;
