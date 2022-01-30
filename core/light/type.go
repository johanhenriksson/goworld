package light

// Type indicates which kind of light. Point, Directional etc
type Type int32

const (
	// AmbientLight is the background light applied to everything.
	Ambient Type = 0

	// PointLight is a normal light casting rays in all directions.
	Point Type = 1

	// DirectionalLight is a directional light source, casting parallell rays.
	Directional Type = 2
)
