package olcode

import (
	"errors"
	"strings"
	"sync"
)

var errInvalidOffset = errors.New("invalid offset")
var errOutOfSzLimit = errors.New("out of the size limit")

const sizeLimit = 4 * 1024 * 1024 // 4MB

// Document is a document
type Document struct {
	ID int64

	mtx     sync.RWMutex
	content string
}

// CheckOffset check whether the offset is valid
func (d *Document) CheckOffset(offset int) error {
	return d.checkOffset(offset, d.Len())
}

func (d *Document) checkOffset(offset, len int) (err error) {
	if offset > len || offset < 0 {
		err = errInvalidOffset
	}
	return
}

func (d *Document) checkSizeLimit(addedSize, len int) (err error) {
	if (len + addedSize) > sizeLimit {
		err = errOutOfSzLimit
	}
	return
}

// Len returns the length of content
func (d *Document) Len() int {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	return len(d.content)
}

// Content returns the current contents of the document
func (d *Document) Content() string {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	return d.content
}

// Insert inserts `text` at the offset of Content
func (d *Document) Insert(offset int, text string) (int, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	cl := len(d.content)
	if err := d.checkOffset(offset, cl); err != nil {
		return 0, err
	}
	if err := d.checkSizeLimit(len(text), cl); err != nil {
		return 0, err
	}
	d.content = strings.Join([]string{d.content[:offset], text, d.content[offset:]}, "")
	return len(text), nil
}

// Delete deletes text of length `size` at the `offset` of the Content
// before controls the deleting direction
func (d *Document) Delete(offset, n int, before bool) (int, int, error) {
	d.mtx.Lock()
	d.mtx.Unlock()

	cl := len(d.content)
	if err := d.checkOffset(offset, cl); err != nil {
		return 0, 0, err
	}
	// [:b] [e:] will be new Content which means [b:e] will be removed
	var b, e int
	switch before {
	case true:
		b = offset - n
		e = offset
	case false:
		b = offset
		e = offset + n
	}

	// fix the problem out of boundary
	if b < 0 {
		b = 0
	}
	if e > cl {
		e = cl
	}
	d.content = strings.Join([]string{d.content[:b], d.content[e:]}, "")
	return b, e, nil
}
