#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/forward_vertex.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)

// Attributes
IN(0, vec3, position)
IN(1, vec3, normal)
IN(2, vec2, texcoord)

void main() 
{
	out_object = gl_InstanceIndex;
	mat4 m = objects.item[out_object].model;

	// texture coords
	out_color.xy = in_texcoord;

	// gbuffer view position
	out_world_position = (m * vec4(in_position.xyz, 1.0)).xyz;
	out_world_position.y += 0.75*sin(out_world_position.x * camera.Time * 0.004);

	out_view_position = (camera.View * vec4(out_world_position, 1)).xyz;

	// world normal
	out_world_normal = normalize((m * vec4(in_normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
