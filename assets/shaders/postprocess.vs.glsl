#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

layout (location = 0) in vec3 inPos;
layout (location = 1) in vec2 texcoord_0;

layout (location = 0) out vec2 texcoord;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	texcoord = texcoord_0;
	gl_Position = vec4(inPos, 1);
}
