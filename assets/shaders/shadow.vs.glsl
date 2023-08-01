#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/material.glsl"

// Attributes
layout (location = 0) in vec3 position;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

layout (location = 0) out float depth;

void main() 
{
	mat4 mvp = camera.ViewProj * objects.item[gl_InstanceIndex].model;
	gl_Position = mvp * vec4(position, 1);

	// store linear depth
	depth = gl_Position.z / gl_Position.w;
}
