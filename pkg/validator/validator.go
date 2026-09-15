package validator

import (
	"regexp"

	"github.com/GiaBao0510/project-templates/pkg/apperror"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail kiểm tra định dạng email.
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return apperror.NewBadRequest("Email không hợp lệ")
	}
	return nil
}

// ValidateRequired kiểm tra field bắt buộc.
func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return apperror.NewBadRequest(fieldName + " không được để trống")
	}
	return nil
}
