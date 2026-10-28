package response

import (
	"encoding/json"
	"net/http"
)

// JSON ghi response theo cùng một quy ước.
//
// Encode trước khi WriteHeader để tránh trường hợp
// đã gửi HTTP 200 nhưng JSON encoding thất bại.
func JSON(
	w http.ResponseWriter,
	status int, 
	data any,
) error {

	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}

	// Dùng Marshal trước khi ghi status vì http.WriteHeader không cho phép đổi status
	payload, err := json.Marshal(data)
	if err != nil{
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(append(payload, '\n'))
	return err
}