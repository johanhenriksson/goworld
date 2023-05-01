#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

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

// including the texture array allows us to use standard materials for the line shader
layout (binding = 2) uniform sampler2D[] Textures;

layout (location = 0) in vec3 position;
layout (location = 1) in vec4 color_0;

layout (location = 0) out vec3 color;

out gl_PerVertex 
{
    vec4 gl_Position;   
};

void main() 
{
    color = color_0.rgb;

	mat4 mvp = camera.ViewProj * ssbo.objects[gl_InstanceIndex].model;
	gl_Position = mvp * vec4(position, 1);
}
