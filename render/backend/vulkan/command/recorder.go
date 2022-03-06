package command

type Recorder interface {
	Record(CommandFn)
	Apply(Buffer)
}

type recorder struct {
	parts []CommandFn
}

func NewRecorder() Recorder {
	return &recorder{
		parts: make([]CommandFn, 0, 16),
	}
}

func (r recorder) Apply(cmd Buffer) {
	for _, part := range r.parts {
		part(cmd)
	}
}

func (r *recorder) Record(cmd CommandFn) {
	r.parts = append(r.parts, cmd)
}
