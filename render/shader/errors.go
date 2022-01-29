package shader

import (
	"errors"
)

var ErrUnknownAttribute = errors.New("unknown attribute")

var ErrUniformType = errors.New("invalid uniform type")
var ErrUnknownUniform = errors.New("unknown uniform")
var ErrUpdateUniform = errors.New("failed to update uniform")
