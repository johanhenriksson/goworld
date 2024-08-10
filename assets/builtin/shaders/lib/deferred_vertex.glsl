// Common vertex shader code
// Varyings
OUT(0, flat uint, object)
OUT(1, vec3, normal)
OUT(2, vec3, position)
OUT(3, vec4, color)

#define OBJECT(idx,name) \
	layout (std430, binding = idx) readonly buffer uniform_ ## name { Object item[]; } _sb_ ## name; \
	Object name = _sb_ ## name.item[gl_InstanceIndex];

out gl_PerVertex 
{
	vec4 gl_Position;   
};
