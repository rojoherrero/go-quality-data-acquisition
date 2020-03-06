package main

type InspectionDocument struct {
	ID            int32
	InspectionID  int64
	Document      []byte
	FileExtension string
}
