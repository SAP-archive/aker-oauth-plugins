package handler

import "fmt"

type InvalidKeyLengthError struct {
	key     []byte
	keyType string
}

func (e *InvalidKeyLengthError) Error() string {
	return fmt.Sprintf("Invalid %s key of length %d, expected %d",
		e.keyType, len(e.key), SessionKeyLength)
}
