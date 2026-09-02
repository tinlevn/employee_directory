package domain

import (
	"github.com/google/uuid"
)

type OrgChartNode struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	JobTitle        *string    `json:"job_title,omitempty"`
	Department      *string    `json:"department,omitempty"`
	ProfilePhotoURL *string    `json:"profile_photo_url,omitempty"`
	ReportsTo       *uuid.UUID `json:"reports_to,omitempty"`
}
