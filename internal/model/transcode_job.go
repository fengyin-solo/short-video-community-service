package model

type TranscodeJob struct {
	ID      string
	Status  string
	Attempt int
	Version int
}
