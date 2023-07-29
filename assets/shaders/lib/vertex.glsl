// Common vertex shader code
// Varyings
layout (location = 0) out vec4 color0;
layout (location = 1) out vec3 normal0;
layout (location = 2) out vec3 position0;
layout (location = 3) out flat uint objectIndex;

out gl_PerVertex 
{
	vec4 gl_Position;   
};
