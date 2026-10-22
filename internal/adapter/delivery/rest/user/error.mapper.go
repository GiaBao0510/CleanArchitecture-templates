package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	shared "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/shared"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
)

func writeError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"

	// Handle specific error types and set appropriate status codes and messages
	switch {
		case errors.Is(err, domain.ErrInvalidName),
			errors.Is(err, domain.ErrInvalidEmail),
			errors.Is(err, domain.ErrInvalidID):
			statusCode = http.StatusBadRequest
			message = err.Error()
		case errors.Is(err, shared.ErrNotFound):
			statusCode = http.StatusNotFound
			message = err.Error()
		case errors.Is(err, domain.ErrEmailIsAlreadyUsed):
			statusCode = http.StatusConflict
			message = err.Error()
		case errors.Is(err, context.Canceled):
			statusCode = 499
			message = "Client Closed Request"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusGatewayTimeout
			message = "Request Timeout"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}