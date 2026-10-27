package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Authorizer interface {
	RequirePermission(permission string) func(next http.Handler) http.Handler
}

func RegisterProtectedUserRoutes(r chi.Router, auth Authorizer) {

	// Đăng ký các route liên quan đến người dùng với middleware xác thực quyền truy cập
	_ = r

	// Sử dụng middleware xác thực quyền truy cập cho các route người dùng
	_ = auth
}
