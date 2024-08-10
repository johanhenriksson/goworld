#version 450

#include "lib/common.glsl"
#include "lib/forward_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)

// Attributes
IN(0, vec3, position)
IN(1, vec3, normal)
IN(2, vec4, color)

void main() 
{
	out_object = gl_InstanceIndex;

	// gbuffer diffuse
	out_color = in_color.rgba;

	// gbuffer view position
	out_view_position = (camera.View * object.model * vec4(in_position.xyz, 1.0)).xyz;
	out_world_position = (object.model * vec4(in_position.xyz, 1.0)).xyz;

	// world normal
	out_world_normal = normalize((object.model * vec4(in_normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
