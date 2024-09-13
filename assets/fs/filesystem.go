package fs

import "fmt"

var ErrNotFound = fmt.Errorf("not found")

type Filesystem interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}
