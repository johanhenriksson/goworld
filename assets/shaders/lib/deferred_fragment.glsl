// Varying
IN(0, flat uint, object)
IN(1, vec3, normal)
IN(2, vec3, position)
IN(3, vec4, color)

// Return Output
OUT(0, vec4, diffuse)
OUT(1, vec4, normal)
OUT(2, vec4, position)

#define OBJECT(idx,name) \
	layout (binding = idx) readonly buffer uniform_ ## name { Object item[]; } __ ## name; \
	Object name = __ ## name.item[in_object];
