package fs

import "fmt"

var ErrNotFound = fmt.Errorf("not found")
var ErrImmutable = fmt.Errorf("immutable")

type Filesystem interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}
