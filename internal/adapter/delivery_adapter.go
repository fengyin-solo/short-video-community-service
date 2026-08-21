package adapter

import (
	"fmt"

	"shortvideo/pkg/delivery"
)

type DeliveryAdapter struct{ client *delivery.Client }

func NewDeliveryAdapter(client *delivery.Client) *DeliveryAdapter {
	return &DeliveryAdapter{client: client}
}

func (a *DeliveryAdapter) Send(videoID string) error {
	if err := a.client.Send(videoID); err != nil {
		return fmt.Errorf("delivery failed: %v", err)
	}
	return nil
}
