package textstore

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases := map[string]struct {
		text       string
		expectedID int64
	}{
		"empty text": {
			text:       "",
			expectedID: 1,
		},
		"some random text": {
			text:       "random string",
			expectedID: 1,
		},
	}
	for name, test := range testCases {
		resetStore()
		t.Run(name, func(t *testing.T) {
			id := Add(test.text)
			if id != test.expectedID {
				t.Fatal("unexpected id")
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	testCases := map[string]struct {
		text        string
		expectedErr string
		skipAdd     bool
	}{
		"empty text": {
			text:        "",
			expectedErr: "",
		},
		"random text": {
			text:        "random string",
			expectedErr: "",
		},
		"item not found": {
			text:        "random string",
			expectedErr: ErrItemNotFound.Error(),
			skipAdd:     true,
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			var id int64
			if !test.skipAdd {
				id = Add(test.text)
			}

			text, err := GetByID(id)
			if len(test.expectedErr) > 0 {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != test.expectedErr {
					t.Fatal("unexpected error message")
				}
				return
			}
			if text != test.text {
				t.Fatal("unexpected text")
			}

		})
	}
}

func resetStore() {
	textStore = sync.Map{}
	id = 0
}
