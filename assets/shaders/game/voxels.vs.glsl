#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/vertex.glsl"

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in uint normal_id;
layout (location = 2) in vec3 color_0;
layout (location = 3) in float occlusion_0;

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
	mat4 mv = camera.View * ssbo.objects[gl_InstanceIndex].model;

	// gbuffer diffuse
	color0 = vec4(color_0, 1 - occlusion_0);

	// gbuffer view space position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	vec3 normal = normals[normal_id];
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
