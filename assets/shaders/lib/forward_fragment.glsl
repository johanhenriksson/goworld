IN(0, flat uint, object)
IN(1, vec4, color)
IN(2, vec3, view_position)
IN(3, vec3, world_normal)
IN(4, vec3, world_position)

// Return Output
OUT(0, vec4, diffuse)

#define OBJECT(idx,name) \
	layout (std430, binding = idx) readonly buffer uniform_ ## name { Object item[]; } __ ## name; \
	Object name = __ ## name.item[in_object];
