package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	usecase "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Struct này đại diện cho một handler trong REST API để xử lý các yêu cầu liên quan đến người dùng. Nó chứa một trường `service` là một con trỏ đến `usecase.Service`, cho phép handler này gọi các phương thức của service để thực hiện các thao tác liên quan đến người dùng, chẳng hạn như tạo, cập nhật hoặc xóa người dùng.
type Handler struct {
	service *usecase.Service
}

// Khởi tạo 
func NewHandler(service *usecase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Hàm decode được sử dụng để giải mã dữ liệu từ yêu cầu HTTP và ánh xạ nó vào một đối tượng đích (dst). 
// Nó nhận vào ba tham số: `w` là ResponseWriter để gửi phản hồi, `r` là Request chứa dữ liệu từ yêu cầu HTTP, 
// và `dst` là đối tượng đích mà dữ liệu sẽ được giải mã vào. Hiện tại, hàm này chưa thực hiện bất kỳ logic nào và chỉ trả về nil, có nghĩa là không có lỗi xảy ra trong quá trình giải mã.
func decode(w http.ResponseWriter, r *http.Request, dst any) error {

	// Bước 1 giới hạn kích thước đầu vào, tối đa là 1MB (1 << 20 bytes) để tránh việc đọc dữ liệu quá lớn từ yêu cầu HTTP
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
	defer r.Body.Close() // Đóng body khi hàm kết thúc

	// Bước 2: Giải mã JSON từ Body
	dec := json.NewDecoder(r.Body)

	// Nếu JSON có field mà struct đích KHÔNG khai báo -> báo lỗi.
	// Ví dụ struct chỉ có "name", "email" mà client gửi thêm "isAdmin": true
	// thì sẽ bị từ chối (mặc định Go sẽ lẳng lặng bỏ qua field thừa).
	dec.DisallowUnknownFields()

	// Bước 3: Giải mã dữ liệu JSON vào đối tượng đích (dst)
	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Bước 4: Kiểm tra xem còn dữ liệu thừa trong body hay không (chỉ nên có 1 object JSON duy nhất)
	var extra any
	err := dec.Decode(&extra)

	// Nếu không còn gì để đọc, Decode trả về io.EOF (End Of File) -> body hợp lệ.
	if !errors.Is(err, io.EOF){
		return errors.New("Request body must only contain a single JSON object")
	}
	
	// Nếu mọi thứ đều hợp lệ, trả về nil (không có lỗi)
	return nil
}

// Hàm respond được sử dụng để gửi phản hồi HTTP với mã trạng thái và dữ liệu JSON.
func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// Hàm parseUUIDFromRequest để trích xuất và phân tích một UUID từ tham số URL của yêu cầu HTTP. Nó sử dụng thư viện chi để lấy giá trị của tham số "id" từ URL và sau đó sử dụng hàm uuid.Parse để chuyển đổi chuỗi thành một đối tượng UUID. Nếu quá trình phân tích thành công, nó trả về UUID và nil (không có lỗi). Nếu có lỗi trong quá trình phân tích, nó trả về một giá trị UUID mặc định và lỗi tương ứng.
func parseUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r,"id"))
}

// Hàm queryInt được sử dụng để trích xuất và phân tích một giá trị số nguyên từ tham số truy vấn của yêu cầu HTTP. Nó lấy giá trị của tham số "limit" từ URL, sau đó sử dụng hàm strconv.Atoi để chuyển đổi chuỗi thành một số nguyên. Nếu quá trình phân tích thành công, nó trả về giá trị số nguyên và nil (không có lỗi). Nếu có lỗi trong quá trình phân tích, nó trả về 0 và lỗi tương ứng.
func queryInt (r *http.Request, key string, fallback int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil || n < 0 || int(64) > 2147483647 {
		return 0, errors.New("Invalid query parameter: " + key)
	}

	return n, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := decode(w, r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.service.Create(r.Context(), req.Input())
	if err != nil {
		writeError(w, err)
		return
	}

	respond(w, http.StatusCreated, toResponse(out))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request){
	id, err := parseUUIDFromRequest(r)

	if err != nil {
		writeError(w, err)
		return
	}

	out, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	respond(w, http.StatusOK, toResponse(out))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, err := queryInt(r, "limit", 20)
	if err != nil || limit <= 0 || limit > 100 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := queryInt(r, "offset", 0)
	if err != nil || offset < 0 {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	out, err := h.service.List(r.Context(), usecase.ListUserRequest{
		Limit: int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	result := make([]UserResponse, len(out))
	for _, u := range out {
		result = append(result, toResponse(u))
	}
	respond(w, http.StatusOK, result)
}
 
func(h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var req UpdateRequest
	if err := decode(w, r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.service.Update(r.Context(), req.Input(id))
	if err != nil {
		writeError(w, err)
		return
	}
	respond(w, http.StatusOK, toResponse(out))
}

func(h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
} 
