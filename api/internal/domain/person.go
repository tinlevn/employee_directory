package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID              uuid.UUID  `json:"id"`
	OrgID           uuid.UUID  `json:"org_id"`
	FirstName       string     `json:"first_name"`
	MiddleName      *string    `json:"middle_name,omitempty"`
	LastName        string     `json:"last_name"`
	PreferredName   *string    `json:"preferred_name,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	ProfilePhotoURL *string    `json:"profile_photo_url,omitempty"`

	PersonalEmail *string `json:"personal_email,omitempty"`
	OrgEmail      *string `json:"org_email,omitempty"`
	PhonePrimary  *string `json:"phone_primary,omitempty"`

	AddressLine1    *string `json:"address_line_1,omitempty"`
	AddressLine2    *string `json:"address_line_2,omitempty"`
	City            *string `json:"city,omitempty"`
	StateProvince   *string `json:"state_province,omitempty"`
	PostalCode      *string `json:"postal_code,omitempty"`
	Country         *string `json:"country,omitempty"`
	IsInternational bool    `json:"is_international"`

	IsActive      bool       `json:"is_active"`
	ArchivedAt    *time.Time `json:"archived_at,omitempty"`
	ArchiveReason *string    `json:"archive_reason,omitempty"`
	Source        *string    `json:"source,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	Tags          []string   `json:"tags"`

	// Enriched from current employment_records (for list view) — not stored in persons table
	CurrentJobTitle       *string    `json:"current_job_title,omitempty"`
	CurrentDepartment     *string    `json:"current_department,omitempty"`
	CurrentTeam           *string    `json:"current_team,omitempty"`
	CurrentOfficeLocation *string    `json:"current_office_location,omitempty"`
	CurrentHireDate       *time.Time `json:"current_hire_date,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

func (p Person) FullName() string {
	if p.MiddleName != nil && *p.MiddleName != "" {
		return p.FirstName + " " + *p.MiddleName + " " + p.LastName
	}
	return p.FirstName + " " + p.LastName
}

// EmergencyContact is a separate entity linked by FK, optional (0..1 per person for now, extensible to many)
type EmergencyContact struct {
	ID           uuid.UUID `json:"id"`
	PersonID     uuid.UUID `json:"person_id"`
	Name         string    `json:"name"`
	Phone        *string   `json:"phone,omitempty"`
	Email        *string   `json:"email,omitempty"`
	Relationship *string   `json:"relationship,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
