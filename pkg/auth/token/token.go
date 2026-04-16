package token

import (
	"encoding"
	"time"
)

type TokenData interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type Token[T TokenData] struct {
	ID       string
	Data     T
	NotAfter time.Time
}

func (t *Token[T]) IsExpired() bool {
	return time.Now().After(t.NotAfter)
}
