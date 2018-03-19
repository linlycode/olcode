package olcode

import (
	"sync"
)

// Document is a document
type Document struct {
	ID int64

	cMtx    sync.Mutex
	Content string
}
