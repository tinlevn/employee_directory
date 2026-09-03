package dto

import "time"

// Pagination
type PaginationQuery struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Sort     string `query:"sort"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}
func (p PaginationQuery) Offset() int { return (p.Page - 1) * p.PageSize }

// Person — trimmed: removed ethnicity/religion/blood_type, maiden_name, pronouns, nationality, phone_secondary, linkedin_url, personal_website_url, inline emergency_contact (now separate entity), cost_center/division/business_unit (employment)
type CreatePersonRequest struct {
	OrgID           string   `json:"org_id" validate:"omitempty,uuid"`
	FirstName       string   `json:"first_name" validate:"required,min=1,max=100"`
	MiddleName      *string  `json:"middle_name" validate:"omitempty,max=100"`
	LastName        string   `json:"last_name" validate:"required,min=1,max=100"`
	PreferredName   *string  `json:"preferred_name" validate:"omitempty,max=100"`
	DateOfBirth     *string  `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Gender          *string  `json:"gender" validate:"omitempty,oneof=male female non-binary prefer-not-to-say"`
	ProfilePhotoURL *string  `json:"profile_photo_url" validate:"omitempty,url,max=2048"`
	PersonalEmail   *string  `json:"personal_email" validate:"omitempty,email,max=255"`
	OrgEmail        *string  `json:"org_email" validate:"omitempty,email,max=255"`
	PhonePrimary    *string  `json:"phone_primary" validate:"omitempty,max=50"`
	AddressLine1    *string  `json:"address_line_1" validate:"omitempty,max=255"`
	AddressLine2    *string  `json:"address_line_2" validate:"omitempty,max=255"`
	City            *string  `json:"city" validate:"omitempty,max=100"`
	StateProvince   *string  `json:"state_province" validate:"omitempty,max=100"`
	PostalCode      *string  `json:"postal_code" validate:"omitempty,max=20"`
	Country         *string  `json:"country" validate:"omitempty,max=100"`
	IsInternational *bool    `json:"is_international"`
	Source          *string  `json:"source" validate:"omitempty,max=100"`
	Notes           *string  `json:"notes" validate:"omitempty,max=5000"`
	Tags            []string `json:"tags" validate:"omitempty,dive,max=50"`
}

type UpdatePersonRequest struct {
	FirstName       *string   `json:"first_name" validate:"omitempty,min=1,max=100"`
	MiddleName      *string   `json:"middle_name" validate:"omitempty,max=100"`
	LastName        *string   `json:"last_name" validate:"omitempty,min=1,max=100"`
	PreferredName   *string   `json:"preferred_name" validate:"omitempty,max=100"`
	DateOfBirth     *string   `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Gender          *string   `json:"gender" validate:"omitempty,oneof=male female non-binary prefer-not-to-say"`
	ProfilePhotoURL *string   `json:"profile_photo_url" validate:"omitempty,url,max=2048"`
	PersonalEmail   *string   `json:"personal_email" validate:"omitempty,email,max=255"`
	OrgEmail        *string   `json:"org_email" validate:"omitempty,email,max=255"`
	PhonePrimary    *string   `json:"phone_primary" validate:"omitempty,max=50"`
	AddressLine1    *string   `json:"address_line_1" validate:"omitempty,max=255"`
	AddressLine2    *string   `json:"address_line_2" validate:"omitempty,max=255"`
	City            *string   `json:"city" validate:"omitempty,max=100"`
	StateProvince   *string   `json:"state_province" validate:"omitempty,max=100"`
	PostalCode      *string   `json:"postal_code" validate:"omitempty,max=20"`
	Country         *string   `json:"country" validate:"omitempty,max=100"`
	IsInternational *bool     `json:"is_international"`
	Source          *string   `json:"source" validate:"omitempty,max=100"`
	Notes           *string   `json:"notes" validate:"omitempty,max=5000"`
	Tags            *[]string `json:"tags" validate:"omitempty"`
}

type ListPersonsQuery struct {
	PaginationQuery
	OrgID      string `query:"org_id" validate:"omitempty,uuid"`
	Search     string `query:"q" validate:"omitempty,max=200"`
	Department string `query:"department" validate:"omitempty,max=200"`
	Location   string `query:"office_location" validate:"omitempty,max=200"`
	Team       string `query:"team" validate:"omitempty,max=200"`
	IsActive   *bool  `query:"is_active"`
	City       string `query:"city" validate:"omitempty,max=100"`
	Country    string `query:"country" validate:"omitempty,max=100"`
	Tag        string `query:"tag" validate:"omitempty,max=50"`
}

// EmergencyContact — separate entity, optional per person
type CreateEmergencyContactRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=200"`
	Phone        *string `json:"phone" validate:"omitempty,max=50"`
	Email        *string `json:"email" validate:"omitempty,email,max=255"`
	Relationship *string `json:"relationship" validate:"omitempty,max=100"`
}

type UpdateEmergencyContactRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=200"`
	Phone        *string `json:"phone" validate:"omitempty,max=50"`
	Email        *string `json:"email" validate:"omitempty,email,max=255"`
	Relationship *string `json:"relationship" validate:"omitempty,max=100"`
}

// Employment — trimmed: removed division, business_unit, cost_center
type CreateEmploymentRequest struct {
	OrgID             string   `json:"org_id" validate:"omitempty,uuid"`
	EmployeeID        *string  `json:"employee_id" validate:"omitempty,max=100"`
	JobTitle          *string  `json:"job_title" validate:"omitempty,max=200"`
	JobLevel          *string  `json:"job_level" validate:"omitempty,max=100"`
	EmploymentStatus  *string  `json:"employment_status" validate:"omitempty,oneof=full-time part-time contract intern freelance probationary on-leave"`
	EmploymentType    *string  `json:"employment_type" validate:"omitempty,oneof=permanent fixed-term casual"`
	WorkArrangement   *string  `json:"work_arrangement" validate:"omitempty,oneof=on-site remote hybrid"`
	Department        *string  `json:"department" validate:"omitempty,max=200"`
	Team              *string  `json:"team" validate:"omitempty,max=200"`
	OfficeLocation    *string  `json:"office_location" validate:"omitempty,max=200"`
	DeskNumber        *string  `json:"desk_number" validate:"omitempty,max=50"`
	ReportsTo         *string  `json:"reports_to" validate:"omitempty,uuid"`
	SalaryAmount      *int64   `json:"salary_amount" validate:"omitempty,min=0"`
	SalaryCurrency    *string  `json:"salary_currency" validate:"omitempty,len=3"`
	PayFrequency      *string  `json:"pay_frequency" validate:"omitempty,oneof=monthly bi-weekly weekly"`
	HourlyRate        *float64 `json:"hourly_rate" validate:"omitempty,min=0"`
	HireDate          *string  `json:"hire_date" validate:"omitempty,datetime=2006-01-02"`
	ProbationEndDate  *string  `json:"probation_end_date" validate:"omitempty,datetime=2006-01-02"`
	ContractStartDate *string  `json:"contract_start_date" validate:"omitempty,datetime=2006-01-02"`
	ContractEndDate   *string  `json:"contract_end_date" validate:"omitempty,datetime=2006-01-02"`
	ValidFrom         string   `json:"valid_from" validate:"required,datetime=2006-01-02"`
}

// Events
type CreateEventRequest struct {
	OrgID          string   `json:"org_id" validate:"omitempty,uuid"`
	EventType      string   `json:"event_type" validate:"required,oneof=HIRED REHIRED PROMOTED DEMOTED TRANSFERRED SECONDMENT_START SECONDMENT_END SALARY_CHANGE TITLE_CHANGE RESIGNED TERMINATED LAID_OFF RETIRED ON_LEAVE_START ON_LEAVE_END CONTRACT_RENEWED CONTRACT_EXPIRED ACTIVATED DEACTIVATED RECORD_UPDATED"`
	Context        *string  `json:"context" validate:"omitempty,oneof=employment general"`
	FromStatus     *string  `json:"from_status" validate:"omitempty,max=100"`
	ToStatus       *string  `json:"to_status" validate:"omitempty,max=100"`
	FromDepartment *string  `json:"from_department" validate:"omitempty,max=200"`
	ToDepartment   *string  `json:"to_department" validate:"omitempty,max=200"`
	FromTitle      *string  `json:"from_title" validate:"omitempty,max=200"`
	ToTitle        *string  `json:"to_title" validate:"omitempty,max=200"`
	FromLocation   *string  `json:"from_location" validate:"omitempty,max=200"`
	ToLocation     *string  `json:"to_location" validate:"omitempty,max=200"`
	Reason         *string  `json:"reason" validate:"omitempty,max=5000"`
	ReasonCode     *string  `json:"reason_code" validate:"omitempty,max=100"`
	IsVoluntary    *bool    `json:"is_voluntary"`
	EffectiveDate  string   `json:"effective_date" validate:"required,datetime=2006-01-02"`
	InitiatedBy    *string  `json:"initiated_by" validate:"omitempty,uuid"`
	ApprovedBy     *string  `json:"approved_by" validate:"omitempty,uuid"`
	WitnessedBy    *string  `json:"witnessed_by" validate:"omitempty,uuid"`
	DocumentURLs   []string `json:"document_urls" validate:"omitempty,dive,url"`
	Notes          *string  `json:"notes" validate:"omitempty,max=5000"`
}

type ListEventsQuery struct {
	PaginationQuery
	EventType string  `query:"event_type" validate:"omitempty,oneof=HIRED REHIRED PROMOTED DEMOTED TRANSFERRED SECONDMENT_START SECONDMENT_END SALARY_CHANGE TITLE_CHANGE RESIGNED TERMINATED LAID_OFF RETIRED ON_LEAVE_START ON_LEAVE_END CONTRACT_RENEWED CONTRACT_EXPIRED ACTIVATED DEACTIVATED RECORD_UPDATED"`
	FromDate  *string `query:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate    *string `query:"to_date" validate:"omitempty,datetime=2006-01-02"`
}

// Transfers
type CreateTransferRequest struct {
	OrgID          string  `json:"org_id" validate:"omitempty,uuid"`
	TransferType   *string `json:"transfer_type" validate:"omitempty,oneof=INTERNAL INTER-CAMPUS SECONDMENT PROMOTION DEMOTION LATERAL"`
	FromDepartment *string `json:"from_department" validate:"omitempty,max=200"`
	ToDepartment   *string `json:"to_department" validate:"omitempty,max=200"`
	FromLocation   *string `json:"from_location" validate:"omitempty,max=200"`
	ToLocation     *string `json:"to_location" validate:"omitempty,max=200"`
	FromManagerID  *string `json:"from_manager_id" validate:"omitempty,uuid"`
	ToManagerID    *string `json:"to_manager_id" validate:"omitempty,uuid"`
	FromTitle      *string `json:"from_title" validate:"omitempty,max=200"`
	ToTitle        *string `json:"to_title" validate:"omitempty,max=200"`
	EffectiveDate  string  `json:"effective_date" validate:"required,datetime=2006-01-02"`
	Reason         *string `json:"reason" validate:"omitempty,max=5000"`
	ApprovedBy     *string `json:"approved_by" validate:"omitempty,uuid"`
	Notes          *string `json:"notes" validate:"omitempty,max=5000"`
}

// Organizations
type CreateOrgRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=255"`
	Type     string  `json:"type" validate:"required,oneof=company ngo government other"`
	Country  *string `json:"country" validate:"omitempty,max=100"`
	Timezone *string `json:"timezone" validate:"omitempty,max=100"`
}

// Analytics
type HeadcountQuery struct {
	OrgID      string  `query:"org_id" validate:"omitempty,uuid"`
	Department *string `query:"department"`
	Location   *string `query:"location"`
	AtDate     *string `query:"at_date" validate:"omitempty,datetime=2006-01-02"`
}

type AttritionQuery struct {
	OrgID    string `query:"org_id" validate:"omitempty,uuid"`
	FromDate string `query:"from_date" validate:"required,datetime=2006-01-02"`
	ToDate   string `query:"to_date" validate:"required,datetime=2006-01-02"`
}

type MovementsQuery struct {
	OrgID    string  `query:"org_id" validate:"omitempty,uuid"`
	FromDate string  `query:"from_date" validate:"required,datetime=2006-01-02"`
	ToDate   string  `query:"to_date" validate:"required,datetime=2006-01-02"`
	GroupBy  *string `query:"group_by" validate:"omitempty,oneof=day month department"`
}

// Helpers to parse dates
func ParseDate(s string) (time.Time, error) { return time.Parse("2006-01-02", s) }
func ParseDatePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
