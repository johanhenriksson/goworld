#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

#include "lib/uniforms.glsl"

// Attributes
layout (location = 0) in vec3 position;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

layout (location = 0) out float depth;

void main() 
{
	mat4 mvp = camera.ViewProj * ssbo.objects[gl_InstanceIndex].model;
	gl_Position = mvp * vec4(position, 1);

	// store linear depth
	depth = gl_Position.z / gl_Position.w;
}
