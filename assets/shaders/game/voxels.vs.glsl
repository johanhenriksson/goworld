#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_vertex.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)

// Attributes
IN(0, vec3, position)
IN(1, uint, normal_id)
IN(2, vec3, color)
IN(3, float, occlusion)

// normal lookup table
const vec3 normals[7] = vec3[7] (
	vec3(0,0,0),  // normal 0 - undefined
	vec3(1,0,0),  // x+
	vec3(-1,0,0), // x-
	vec3(0,1,0),  // y+
	vec3(0,-1,0), // y-
	vec3(0,0,1),  // z+
	vec3(0,0,-1)  // z-
);

void main() 
{
	mat4 mv = camera.View * objects.item[gl_InstanceIndex].model;

	// gbuffer diffuse
	out_color = vec4(in_color, 1 - in_occlusion);

	// gbuffer view space position
	out_position = (mv * vec4(in_position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	vec3 normal = normals[in_normal_id];
	out_normal = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
