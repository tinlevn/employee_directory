package dto

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ErrorResponse struct {
	Type    string            `json:"type"`
	Title   string            `json:"title"`
	Status  int               `json:"status"`
	Detail  string            `json:"detail,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	TraceID string            `json:"trace_id,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type HeadcountResponse struct {
	Department string `json:"department"`
	Count      int    `json:"count"`
}

type DimensionCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type AttritionResponse struct {
	PeriodStart      string  `json:"period_start"`
	PeriodEnd        string  `json:"period_end"`
	TotalExits       int64   `json:"total_exits"`
	VoluntaryExits   int64   `json:"voluntary_exits"`
	InvoluntaryExits int64   `json:"involuntary_exits"`
	AvgHeadcount     float64 `json:"avg_headcount"`
	AttritionRatePct float64 `json:"attrition_rate_pct"`
}

type MovementPoint struct {
	Date       string `json:"date"`
	NewEntries int    `json:"new_entries"`
	Exits      int    `json:"exits"`
	NetChange  int    `json:"net_change"`
}

type SnapshotResponse struct {
	Date       string              `json:"date"`
	OrgID      string              `json:"org_id"`
	Headcount  int                 `json:"headcount"`
	ByDept     []HeadcountResponse `json:"by_department"`
	ByLocation []DimensionCount    `json:"by_location"`
}
