package textstore

import (
	"errors"
	"sync"
	"sync/atomic"
)

var id int64
var textStore sync.Map

var ErrItemNotFound = errors.New("item not found from store")

func Add(text string) int64 {
	id := atomic.AddInt64(&id, 1)
	textStore.Store(id, text)

	return id
}

func GetByID(id int64) (string, error) {
	typelessText, ok := textStore.Load(id)
	if !ok {
		return "", ErrItemNotFound
	}

	textStr, _ := typelessText.(string)

	return textStr, nil
}
