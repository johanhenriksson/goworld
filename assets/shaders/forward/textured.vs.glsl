#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/vertex.glsl"

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texcoord_0;

// Varyings
layout (location = 4) out vec3 wnormal;
layout (location = 5) out vec3 wposition;

void main() 
{
	mat4 m = objects.item[gl_InstanceIndex].model;
	mat4 mv = camera.View * m;

	// texture coords
	color0.xy = texcoord_0;

	// gbuffer view position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;
	wposition = (m * vec4(position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// world normal
	wnormal = normalize((m * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
