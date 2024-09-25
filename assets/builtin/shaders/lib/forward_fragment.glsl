IN(0, flat uint, object)
IN(1, vec4, color)
IN(2, vec2, texcoord)
IN(3, vec3, view_position)
IN(4, vec3, world_normal)
IN(5, vec3, world_position)

// Return Output
OUT(0, vec4, diffuse)

#define OBJECT(idx,name) \
	layout (scalar, binding = idx) readonly buffer uniform_ ## name { Object item[]; } _sb_ ## name; \
	Object name = _sb_ ## name.item[in_object];

#define TEX_SLOT_DIFFUSE 0
