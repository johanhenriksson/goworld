#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"
#include "lib/deferred_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)

// Attributes
IN(0, vec3, position)
IN(1, float, tex_x)
IN(2, vec3, normal)
IN(3, float, tex_y)
IN(4, vec4, color)

void main() 
{
	out_object = gl_InstanceIndex;
	mat4 mv = camera.View * object.model;

	// textures
	out_color.xy = vec2(in_tex_x, in_tex_y);

	// gbuffer position
	out_position = (mv * vec4(in_position, 1)).xyz;

	// gbuffer normal
	out_normal = normalize((mv * vec4(in_normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
