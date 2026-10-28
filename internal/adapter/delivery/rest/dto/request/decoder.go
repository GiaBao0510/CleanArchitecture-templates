package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const MaxBodySize int64 = 1 << 20 // 1Mb

// DecodeJSON giải mã đúng một JSON object.
// T là kiểu request DTO của từng feature.
func DecodeJSON[T any](
	w http.ResponseWriter,
	r *http.Request,
) (T, error) {

	var dst T

	// Bước 1 giới hạn kích thước đầu vào, tối đa là 1MB (1 << 20 bytes) để tránh việc đọc dữ liệu quá lớn từ yêu cầu HTTP
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	defer r.Body.Close()

	// Bước 2: Giải mã JSON từ Body
	decoder := json.NewDecoder(r.Body)

	// Nếu JSON có field mà struct đích KHÔNG khai báo -> báo lỗi.
	// Ví dụ struct chỉ có "name", "email" mà client gửi thêm "isAdmin": true
	// thì sẽ bị từ chối (mặc định Go sẽ lẳng lặng bỏ qua field thừa).
	decoder.DisallowUnknownFields()

	// Bước 3: Giải mã dữ liệu JSON vào đối tượng đích (dst)
	if err := decoder.Decode(&dst); err != nil {
		return dst, fmt.Errorf("decode JSON: %w", err)
	}

	// Bước 4: Kiểm tra xem còn dữ liệu thừa trong body hay không (chỉ nên có 1 object JSON duy nhất)
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) { // Nếu không còn gì để đọc, Decode trả về io.EOF (End Of File) -> body hợp lệ.
		return dst, errors.New("request must contain one JSON value")
	}

	// Nếu mọi thứ đều hợp lệ, trả về nil (không có lỗi)
	return dst, nil
}