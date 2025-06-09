package math_rand

import "github.com/google/uuid"

var RandomUuid = func() string {
	return uuid.NewString()
}
