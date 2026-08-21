package model

type VideoDeliveryJob struct {
	Region    string
	Duplicate bool
	Release   <-chan struct{}
}

type VideoDeliveryResult struct {
	Region string
}
