#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Attributes
layout (location = 0) in vec3 position;

// Uniforms
layout (binding = 0) uniform Camera {
	mat4 Proj;
	mat4 View;
	mat4 ViewProj;
	mat4 ProjInv;
	mat4 ViewInv;
	mat4 ViewProjInv;
	vec3 Eye;
	vec3 Forward;
} camera;

struct ObjectData{
	mat4 model;
	uint textures[4];
};

layout (binding = 1) readonly buffer ObjectBuffer {
	ObjectData objects[];
} ssbo;

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
