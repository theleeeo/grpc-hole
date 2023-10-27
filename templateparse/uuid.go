package templateparse

import (
	"github.com/google/uuid"
)

func uuidFunc() string {
	return uuid.NewString()
}
