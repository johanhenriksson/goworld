#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/vertex.glsl"

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texcoord_0;

void main() 
{
	objectIndex = gl_InstanceIndex;
	mat4 mv = camera.View * objects.item[objectIndex].model;

	// textures
	color0.xy = texcoord_0;

	// gbuffer position
	position0 = (mv * vec4(position, 1)).xyz;

	// gbuffer normal
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
