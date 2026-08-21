package model

// PreviewRequest describes one asynchronous video preview preparation.
type PreviewRequest struct {
	VideoID string
	DelayMS int
}

// PreviewResult is returned after preview preparation finishes.
type PreviewResult struct {
	VideoID string
	Ready   bool
}
