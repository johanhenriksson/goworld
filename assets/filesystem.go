package assets

type Filesystem interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}
