package model

import "time"

type OracleCloudResource struct {
	Identifier   string            `json:"Identifier,omitempty"`
	ResourceType string            `json:"ResourceType,omitempty"`
	DisplayName  string            `json:"DisplayName,omitempty"`
	TimeCreated  *time.Time        `json:"TimeCreated"`
	DefinedTags  map[string]string `json:"DefinedTags,omitempty"`
}
