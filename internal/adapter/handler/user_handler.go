package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/GiaBao0510/project-templates/internal/adapter/handler/dto"
	"github.com/GiaBao0510/project-templates/internal/usecase"
	"github.com/GiaBao0510/project-templates/pkg/apperror"
	"github.com/GiaBao0510/project-templates/pkg/response"
)

// UserHandler xử lý các HTTP request liên quan đến User.
// Handler thuộc tầng adapter — chuyển đổi HTTP request thành lời gọi use case.
type UserHandler struct {
	userUC *usecase.UserUsecase
}

// NewUserHandler khởi tạo handler với dependency injection.
func NewUserHandler(userUC *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// Create godoc
// @Summary Tạo user mới
// @Tags    users
// @Accept  json
// @Produce json
// @Param   body body dto.CreateUserRequest true "User data"
// @Router  /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", err.Error())
		return
	}

	user, err := h.userUC.Create(c.Request.Context(), req.FullName, req.Email)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo user thành công", dto.ToUserResponse(user))
}

// GetByID godoc
// @Summary Lấy user theo ID
// @Tags    users
// @Produce json
// @Param   id path string true "User ID"
// @Router  /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID không hợp lệ", err.Error())
		return
	}

	user, err := h.userUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Thành công", dto.ToUserResponse(user))
}

// GetAll godoc
// @Summary Lấy danh sách users
// @Tags    users
// @Produce json
// @Router  /api/v1/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userUC.GetAll(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Thành công", dto.ToUserResponseList(users))
}

// Update godoc
// @Summary Cập nhật user
// @Tags    users
// @Accept  json
// @Produce json
// @Param   id   path string              true "User ID"
// @Param   body body dto.UpdateUserRequest true "User data"
// @Router  /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID không hợp lệ", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", err.Error())
		return
	}

	user, err := h.userUC.Update(c.Request.Context(), id, req.FullName, req.Email)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Cập nhật thành công", dto.ToUserResponse(user))
}

// Delete godoc
// @Summary Xóa user
// @Tags    users
// @Produce json
// @Param   id path string true "User ID"
// @Router  /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID không hợp lệ", err.Error())
		return
	}

	if err := h.userUC.Delete(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Xóa thành công", nil)
}
