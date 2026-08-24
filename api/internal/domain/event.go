package domain

import (
	"time"

	"github.com/google/uuid"
)

type StatusChangeEvent struct {
	ID               uuid.UUID  `json:"id"`
	PersonID         uuid.UUID  `json:"person_id"`
	OrgID            uuid.UUID  `json:"org_id"`
	EventType        string     `json:"event_type"`
	Context          string     `json:"context"`
	FromStatus       *string    `json:"from_status,omitempty"`
	ToStatus         *string    `json:"to_status,omitempty"`
	FromDepartment   *string    `json:"from_department,omitempty"`
	ToDepartment     *string    `json:"to_department,omitempty"`
	FromTitle        *string    `json:"from_title,omitempty"`
	ToTitle          *string    `json:"to_title,omitempty"`
	FromLocation     *string    `json:"from_location,omitempty"`
	ToLocation       *string    `json:"to_location,omitempty"`
	Reason           *string    `json:"reason,omitempty"`
	ReasonCode       *string    `json:"reason_code,omitempty"`
	IsVoluntary      *bool      `json:"is_voluntary,omitempty"`
	EffectiveDate    time.Time  `json:"effective_date"`
	RecordedAt       time.Time  `json:"recorded_at"`
	InitiatedBy      *uuid.UUID `json:"initiated_by,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	WitnessedBy      *uuid.UUID `json:"witnessed_by,omitempty"`
	LinkedRecordID   *uuid.UUID `json:"linked_record_id,omitempty"`
	LinkedRecordType *string    `json:"linked_record_type,omitempty"`
	DocumentURLs     []string   `json:"document_urls"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
}

type TransferRecord struct {
	ID             uuid.UUID  `json:"id"`
	PersonID       uuid.UUID  `json:"person_id"`
	OrgID          uuid.UUID  `json:"org_id"`
	TransferType   *string    `json:"transfer_type,omitempty"`
	FromDepartment *string    `json:"from_department,omitempty"`
	ToDepartment   *string    `json:"to_department,omitempty"`
	FromLocation   *string    `json:"from_location,omitempty"`
	ToLocation     *string    `json:"to_location,omitempty"`
	FromManagerID  *uuid.UUID `json:"from_manager_id,omitempty"`
	ToManagerID    *uuid.UUID `json:"to_manager_id,omitempty"`
	FromTitle      *string    `json:"from_title,omitempty"`
	ToTitle        *string    `json:"to_title,omitempty"`
	EffectiveDate  time.Time  `json:"effective_date"`
	Reason         *string    `json:"reason,omitempty"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
}

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Country   *string   `json:"country,omitempty"`
	Timezone  *string   `json:"timezone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HeadcountSnapshot struct {
	ID           uuid.UUID `json:"id"`
	OrgID        uuid.UUID `json:"org_id"`
	SnapshotDate time.Time `json:"snapshot_date"`
	Department   *string   `json:"department,omitempty"`
	Location     *string   `json:"location,omitempty"`
	ActiveCount  int       `json:"active_count"`
	NewEntries   int       `json:"new_entries"`
	Exits        int       `json:"exits"`
	NetChange    int       `json:"net_change"`
	CreatedAt    time.Time `json:"created_at"`
}
