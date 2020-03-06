package main

// FailureGroups
type FailureGroups struct {
	Code    string
	Name    string
	Comment string
}

// Failure
type Failure struct {
	Code             string
	FailureGroupCode string
	Nname            string
	Description      string
}
