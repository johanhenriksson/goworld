#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

// Varying
IN(0, flat uint, object)
IN(1, vec3, normal)
IN(2, vec3, position)
IN(3, vec4, color)

// Return Output
OUT(0, vec4, diffuse)
OUT(1, vec4, normal)
OUT(2, vec4, position)

OBJECT(1, object, in_object)

void main() 
{
	out_diffuse = in_color;
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
