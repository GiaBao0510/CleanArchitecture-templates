package errors

import (
	"encoding/json"
	"net/http"
	appError "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/shared/apperror"
)

type ErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, err error) {
	status := Status(err)

	message := "internal server error"
	code := "INTERNAL"

	if status < 500 {
		message = err.Error()
		code = string(appError.KindOf(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code: code,
		Message: message,
	})
}