package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/segmentio/kafka-go"
)

type Publisher struct{
	writer *kafka.Writer
}

// Hàm khởi tạo
func NewPublisher(w *kafka.Writer) *Publisher {
	return &Publisher{writer: w}
}

// Hàm đẩy message lên hàng đợi
func(p *Publisher) Publish(
	ctx context.Context,
	event domain.Event,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal domain event: %w", err)
	}

	// Cấu hình thông tin message
	msg := kafka.Message{
		Key: event.AggegateID[:],
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "event_ID", Value: event.ID[:]},
		},
	}

	// Thực hiện gửi
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("publish %s: %w", event.Type, err)
	}

	return nil
}

