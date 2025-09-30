package models

type MedicineRecord struct {
	TM2Code         string  `json:"tm2_code" csv:"tm2_code"`
	Code            string  `json:"code" csv:"code"`
	TM2Title        string  `json:"tm2_title" csv:"tm2_title"`
	TM2Definition   string  `json:"tm2_definition" csv:"tm2_definition"`
	CodeTitle       string  `json:"code_title" csv:"code_title"`
	Description     string  `json:"code_description" csv:"code_description"`
	ConfidenceScore float64 `json:"confidence_score" csv:"confidence_score"`
	Type            string  `json:"type" csv:"type"`
	TM2Link         string  `json:"tm2_link" csv:"tm2_link"`
}

// Keep existing for compatibility if needed
type FHIRParameters struct {
	ResourceType string      `json:"resourceType"`
	ID           string      `json:"id"`
	Meta         FHIRMeta    `json:"meta"`
	Parameter    []Parameter `json:"parameter"`
}

type FHIRMeta struct {
	VersionId   string `json:"versionId"`
	LastUpdated string `json:"lastUpdated"`
}

type Parameter struct {
	Name         string      `json:"name"`
	ValueBoolean *bool       `json:"valueBoolean,omitempty"`
	ValueInteger *int        `json:"valueInteger,omitempty"`
	ValueString  *string     `json:"valueString,omitempty"`
	ValueCode    *string     `json:"valueCode,omitempty"`
	ValueDecimal *float64    `json:"valueDecimal,omitempty"`
	Part         []Parameter `json:"part,omitempty"`
}