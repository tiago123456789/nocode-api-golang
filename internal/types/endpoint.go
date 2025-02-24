package types

type ValidationsToApply struct {
	Field     string `json:"field"`
	TypeField string `json:"typeField"`
	Rules     string `json:"rules"`
}

type ActionBeforePersist struct {
	Field  string `json:"field"`
	Action string `json:"action"`
}

type Endpoint struct {
	ID                   int64                 `json:"id"`
	Table                string                `json:"table"`
	Path                 string                `json:"path"`
	IsPublic             bool                  `json:"isPublic"`
	Query                string                `json:"query"`
	QueryParams          []string              `json:"queryParams"`
	Validations          []ValidationsToApply  `json:"validations"`
	ActionsBeforePersist []ActionBeforePersist `json:"actionsBeforePersist"`
	IsCacheable          bool                  `json:"isCacheable"`
	CacheTtl             int                   `json:"cacheTtl"`
}
