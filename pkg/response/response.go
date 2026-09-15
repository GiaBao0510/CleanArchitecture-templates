package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GiaBao0510/project-templates/pkg/apperror"
)

// Response là cấu trúc JSON response chuẩn cho toàn bộ API.
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success trả về response thành công.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Error trả về response lỗi.
func Error(c *gin.Context, status int, message string, detail interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Error:   detail,
	})
}

// HandleError xử lý error và trả về response phù hợp.
// Tự động nhận diện AppError để lấy đúng HTTP status code.
func HandleError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Message, nil)
		return
	}
	Error(c, http.StatusInternalServerError, "Lỗi hệ thống", nil)
}
