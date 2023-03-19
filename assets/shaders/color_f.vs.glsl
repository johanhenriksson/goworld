#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec4 color_0;

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
};

layout (binding = 1) readonly buffer ObjectBuffer {
	ObjectData objects[];
} ssbo;

layout (binding = 2) uniform sampler2D[] Textures;

// Varyings
layout (location = 0) out vec4 color0;
layout (location = 1) out vec3 normal0;
layout (location = 2) out vec3 position0;
layout (location = 3) out vec3 wnormal;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	mat4 m = ssbo.objects[gl_InstanceIndex].model;
	mat4 mv = camera.View * m;

	// gbuffer diffuse
	color0 = color_0.rgba;

	// gbuffer position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// world normal
	wnormal = normalize((m * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}