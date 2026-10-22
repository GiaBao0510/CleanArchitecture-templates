package shared

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type ID string

// NewID generates a new unique ID for a user.
func NewID() ( ID, error) {
	var b [16]byte

	// If there is an error generating random bytes, return an empty ID and the error.
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return ID(fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	)), nil
}



// Hàm ParseID nhận một chuỗi và kiểm tra xem nó có phải là định dạng ID hợp lệ hay không. Nếu hợp lệ, nó trả về ID dưới dạng kiểu ID, nếu không hợp lệ, nó trả về lỗi ErrInvalidID.
const hexChars = "0123456789abcdef"
func ParseID(value string) (ID, error) {
	if value == ""  || len(value) != 36  {
		return "", ErrInvalidID
	}
	
	value = strings.ToLower(value)
	for i, c := range value {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return "", ErrInvalidID
			}
		default:
			if !strings.ContainsRune(hexChars, c) {
				return "", ErrInvalidID
			}
		}
	}

	return ID(value), nil
}