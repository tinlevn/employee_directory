package domain

// Gender
type Gender string

const (
	GenderMale         Gender = "male"
	GenderFemale       Gender = "female"
	GenderNonBinary    Gender = "non-binary"
	GenderPreferNotSay Gender = "prefer-not-to-say"
)

func (g Gender) Valid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderNonBinary, GenderPreferNotSay:
		return true
	}
	return false
}

// Employment
type EmploymentStatus string

const (
	EmpFullTime     EmploymentStatus = "full-time"
	EmpPartTime     EmploymentStatus = "part-time"
	EmpContract     EmploymentStatus = "contract"
	EmpIntern       EmploymentStatus = "intern"
	EmpFreelance    EmploymentStatus = "freelance"
	EmpProbationary EmploymentStatus = "probationary"
	EmpOnLeave      EmploymentStatus = "on-leave"
)

type EmploymentType string

const (
	EmploymentPermanent EmploymentType = "permanent"
	EmploymentFixedTerm EmploymentType = "fixed-term"
	EmploymentCasual    EmploymentType = "casual"
)

type WorkArrangement string

const (
	WorkOnSite WorkArrangement = "on-site"
	WorkRemote WorkArrangement = "remote"
	WorkHybrid WorkArrangement = "hybrid"
)

type PayFrequency string

const (
	PayMonthly  PayFrequency = "monthly"
	PayBiWeekly PayFrequency = "bi-weekly"
	PayWeekly   PayFrequency = "weekly"
)

// Events
type EventType string

const (
	EventHired           EventType = "HIRED"
	EventRehired         EventType = "REHIRED"
	EventPromoted        EventType = "PROMOTED"
	EventDemoted         EventType = "DEMOTED"
	EventTransferred     EventType = "TRANSFERRED"
	EventSecondmentStart EventType = "SECONDMENT_START"
	EventSecondmentEnd   EventType = "SECONDMENT_END"
	EventSalaryChange    EventType = "SALARY_CHANGE"
	EventTitleChange     EventType = "TITLE_CHANGE"
	EventResigned        EventType = "RESIGNED"
	EventTerminated      EventType = "TERMINATED"
	EventLaidOff         EventType = "LAID_OFF"
	EventRetired         EventType = "RETIRED"
	EventOnLeaveStart    EventType = "ON_LEAVE_START"
	EventOnLeaveEnd      EventType = "ON_LEAVE_END"
	EventContractRenewed EventType = "CONTRACT_RENEWED"
	EventContractExpired EventType = "CONTRACT_EXPIRED"
	EventActivated       EventType = "ACTIVATED"
	EventDeactivated     EventType = "DEACTIVATED"
	EventRecordUpdated   EventType = "RECORD_UPDATED"
)

var AllEventTypes = []EventType{
	EventHired, EventRehired, EventPromoted, EventDemoted, EventTransferred,
	EventSecondmentStart, EventSecondmentEnd, EventSalaryChange, EventTitleChange,
	EventResigned, EventTerminated, EventLaidOff, EventRetired,
	EventOnLeaveStart, EventOnLeaveEnd, EventContractRenewed, EventContractExpired,
	EventActivated, EventDeactivated, EventRecordUpdated,
}

func (e EventType) Valid() bool {
	switch e {
	case EventHired, EventRehired, EventPromoted, EventDemoted, EventTransferred,
		EventSecondmentStart, EventSecondmentEnd, EventSalaryChange, EventTitleChange,
		EventResigned, EventTerminated, EventLaidOff, EventRetired,
		EventOnLeaveStart, EventOnLeaveEnd, EventContractRenewed, EventContractExpired,
		EventActivated, EventDeactivated, EventRecordUpdated:
		return true
	}
	return false
}

type EventContext string

const (
	ContextEmployment EventContext = "employment"
	ContextGeneral    EventContext = "general"
)

// Transfers
type TransferType string

const (
	TransferInternal    TransferType = "INTERNAL"
	TransferInterCampus TransferType = "INTER-CAMPUS"
	TransferSecondment  TransferType = "SECONDMENT"
	TransferPromotion   TransferType = "PROMOTION"
	TransferDemotion    TransferType = "DEMOTION"
	TransferLateral     TransferType = "LATERAL"
)

// Org type
type OrgType string

const (
	OrgCompany    OrgType = "company"
	OrgNGO        OrgType = "ngo"
	OrgGovernment OrgType = "government"
	OrgOther      OrgType = "other"
)

// Role
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleStaff    Role = "staff"
	RoleReadOnly Role = "read-only"
)
