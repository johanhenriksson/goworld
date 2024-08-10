#version 450

#include "lib/common.glsl"
#include "lib/deferred_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)

// Attributes
IN(0, vec3, position)
IN(1, uint, texdata)
IN(2, vec3, color)
IN(3, float, occlusion)

OUT(4, vec2, texcoord0)
OUT(5, flat uint, texture_id)

// normal lookup table
const uint texdata_normal_mask = 0x07;
const vec3 normals[7] = vec3[7] (
	vec3(0,0,0),  // normal 0 - undefined
	vec3(1,0,0),  // x+
	vec3(-1,0,0), // x-
	vec3(0,1,0),  // y+
	vec3(0,-1,0), // y-
	vec3(0,0,1),  // z+
	vec3(0,0,-1)  // z-
);

const uint texdata_uv_mask = 0x18;
const vec2 uvs[4] = vec2[4] (
	vec2(0, 0), // top left
	vec2(0, 1), // top right
	vec2(1, 0), // bottom left
	vec2(1, 1)  // bottom right
);

const uint texdata_texture_mask = 0xE0;

void main() 
{
	out_object = gl_InstanceIndex;
	mat4 mv = camera.View * object.model;

	// gbuffer diffuse
	out_color = vec4(in_color, 1 - in_occlusion);

	// gbuffer view space position
	out_position = (mv * vec4(in_position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	uint in_normal_id = in_texdata & texdata_normal_mask;
	vec3 normal = normals[in_normal_id];
	out_normal = normalize((mv * vec4(normal, 0.0)).xyz);

	out_texcoord0 = uvs[(in_texdata & texdata_uv_mask) >> 3];
	out_texture_id = (in_texdata & texdata_texture_mask) >> 5;

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
