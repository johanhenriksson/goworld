package light

// Type indicates which kind of light. Point, Directional etc
type Type int32

const (
	// PointLight is a normal light casting rays in all directions.
	TypePoint Type = 1

	// DirectionalLight is a directional light source, casting parallell rays.
	TypeDirectional Type = 2
)
