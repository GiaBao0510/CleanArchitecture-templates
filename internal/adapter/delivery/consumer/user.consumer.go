package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Consumer lắng nghe và xử lý message
type UserEventHandler interface {
	HandlerUserEvent(context.Context, domain.Event) error
}

type UserConsumer struct {
	reader  *kafka.Reader
	handler UserEventHandler
	logger  *zap.Logger
}

func NewUserConsumer(
	reader *kafka.Reader,
	handler UserEventHandler,
	logger *zap.Logger,
) *UserConsumer {
	return &UserConsumer{
		handler: handler,
		reader:  reader,
		logger:  logger,
	}
}

// Hàm Run: chủ yếu để phát hiện vấn đề nếu gặp lỗi trong quá trình thực hiện xử lý message sẽ thông báo lỗi
func (c *UserConsumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("fetch Kafka message: %w", err)
		}

		// Xử lý message
		if err := c.handle(ctx, msg); err != nil {
			c.logger.Error("process user event failed",
				zap.Error(err),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
			)
			// Do not commit a failed message.
			// Stop the worker to avoid silently skipping it.
			return err
		}

		// Thực hiện commit sau khi xử lý xong message
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit Kafka offset: %w", err)
		}
	}
}

// Hàm xử lý message
func (c *UserConsumer) handle(
	ctx context.Context,
	msg kafka.Message,
) error {
	var event domain.Event

	// Chuyển đổi message sang dạng khớp với struct
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("decode event: %w", err)
	}

	// dự vào dạng loại evnet để xửu lý cho chính xác
	switch event.Type {
	case domain.EventTypeUserCreated,
		domain.EventTypeUserUpdated,
		domain.EventTypeUserDeleted:
		return c.handler.HandlerUserEvent(ctx, event)
	
	default:// Event không tồn tại
		return fmt.Errorf("unsupported event type: %s", event.Type)
	}
}

// Hàm đóng
func (c *UserConsumer) Close() error {
	return c.reader.Close()
}
