// Common vertex shader code
// Varyings
OUT(0, flat uint, object)
OUT(1, vec3, normal)
OUT(2, vec3, position)
OUT(3, vec4, color)

out gl_PerVertex 
{
	vec4 gl_Position;   
};
