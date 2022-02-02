package shader

type StageID uint32

type Stage interface {
	ID() StageID
	Compile(source, path string) error
	CompileFile(path string) error
}
