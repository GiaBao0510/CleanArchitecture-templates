package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Khởi tạo một logger mới dựa trên môi trường (production hoặc development).
func New(environment string) (*zap.Logger, error){
	if environment == "production" {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create production logger: %w", err)
		}

		return logger, nil
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create development logger: %w", err)
	}
	return logger, nil
}