#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_GOOGLE_include_directive : enable

#include "lib/uniforms.glsl"

// Varying
layout (location = 0) in vec3 color0;
layout (location = 1) in vec3 normal0;
layout (location = 2) in vec3 position0;

// Return Output
layout (location = 0) out vec4 diffuse;
layout (location = 1) out vec4 normal;
layout (location = 2) out vec4 position;

void main() 
{
	diffuse = vec4(color0, 1);

	vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
	normal = pack_normal;

	position = vec4(position0, 1);
}
