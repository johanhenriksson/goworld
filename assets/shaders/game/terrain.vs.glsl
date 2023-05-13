#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texcoord_0;

// Uniforms
layout (binding = 0) uniform Camera {
	mat4 Proj;
	mat4 View;
	mat4 ViewProj;
	mat4 ProjInv;
	mat4 ViewInv;
	mat4 ViewProjInv;
	vec3 Eye;
} camera;

struct ObjectData{
	mat4 model;
	uint textures[4];
};

layout (binding = 1) readonly buffer ObjectBuffer {
	ObjectData objects[];
} ssbo;

layout (binding = 2) uniform sampler2D[] Textures;

// Varyings
layout (location = 0) out flat uint objectIndex;
layout (location = 1) out vec3 normal0;
layout (location = 2) out vec3 position0;
layout (location = 3) out vec2 texcoord0;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	objectIndex = gl_InstanceIndex;
	mat4 mv = camera.View * ssbo.objects[objectIndex].model;

	// textures
	texcoord0 = texcoord_0;
	uint texture0 = ssbo.objects[objectIndex].textures[0];

	// gbuffer position
	position0 = (mv * vec4(position, 1)).xyz;

	// gbuffer normal
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}