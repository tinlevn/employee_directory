package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentRecord struct {
	ID                uuid.UUID  `json:"id"`
	PersonID          uuid.UUID  `json:"person_id"`
	OrgID             uuid.UUID  `json:"org_id"`
	EmployeeID        *string    `json:"employee_id,omitempty"`
	JobTitle          *string    `json:"job_title,omitempty"`
	JobLevel          *string    `json:"job_level,omitempty"`
	EmploymentStatus  *string    `json:"employment_status,omitempty"`
	EmploymentType    *string    `json:"employment_type,omitempty"`
	WorkArrangement   *string    `json:"work_arrangement,omitempty"`
	Department        *string    `json:"department,omitempty"`
	Team              *string    `json:"team,omitempty"`
	OfficeLocation    *string    `json:"office_location,omitempty"`
	DeskNumber        *string    `json:"desk_number,omitempty"`
	ReportsTo         *uuid.UUID `json:"reports_to,omitempty"`
	SalaryAmount      *int64     `json:"salary_amount,omitempty"`
	SalaryCurrency    *string    `json:"salary_currency,omitempty"`
	PayFrequency      *string    `json:"pay_frequency,omitempty"`
	HourlyRate        *float64   `json:"hourly_rate,omitempty"`
	HireDate          *time.Time `json:"hire_date,omitempty"`
	ProbationEndDate  *time.Time `json:"probation_end_date,omitempty"`
	ContractStartDate *time.Time `json:"contract_start_date,omitempty"`
	ContractEndDate   *time.Time `json:"contract_end_date,omitempty"`
	ValidFrom         time.Time  `json:"valid_from"`
	ValidTo           *time.Time `json:"valid_to,omitempty"`
	IsCurrent         bool       `json:"is_current"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}
