OUT(0, flat uint, object)
OUT(1, vec4, color)
OUT(2, vec3, view_position)
OUT(3, vec3, world_normal)
OUT(4, vec3, world_position)

#define OBJECT(idx,name) \
	layout (std430, binding = idx) readonly buffer uniform_ ## name { Object item[]; } _sb_ ## name; \
	Object name = _sb_ ## name.item[gl_InstanceIndex];

out gl_PerVertex 
{
	vec4 gl_Position;   
};
