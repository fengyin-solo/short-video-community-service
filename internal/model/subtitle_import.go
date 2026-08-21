package model

type SubtitleFile struct {
	Name    string
	Content string
}

type SubtitleImportAudit struct {
	BatchID string
	Status  string
	Failure string
}
