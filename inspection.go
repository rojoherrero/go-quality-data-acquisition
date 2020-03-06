package main

import (
	"time"
)

// Inspection
type Inspection struct {
	ID                int64
	ProductionOrderID string
	PartID            int32
	InspectorID       int64
	InspectorName     string
	FailureCode       string
	Passed            bool
	Comment           string
	DateTime          time.Time
}

// InspectionDTO
type InspectionDTO struct {
	ProductionOrderID string    `json:"production_order_id,omitempty"`
	PartID            string    `json:"part_id,omitempty"`
	InspectionID      int32     `json:"inspection_id,omitempty"`
	Timestamp         time.Time `json:"timestamp,omitempty"`
	InspectorID       int32     `json:"inspector_id,omitempty"`
	InspectorName     string    `json:"inspector_name,omitempty"`
	Errors            []string  `json:"errors,omitempty"`
	Comments          string    `json:"comments,omitempty"`
	Failed            bool      `json:"failed,omitempty"`
}
