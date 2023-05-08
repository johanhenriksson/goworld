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
layout (location = 0) out float depth;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	uint objectIndex = gl_InstanceIndex;
	uint texture0 = ssbo.objects[objectIndex].textures[0];

	mat4 mvp = camera.ViewProj * ssbo.objects[objectIndex].model;

	// gbuffer position
	float center = texture(Textures[texture0], texcoord_0).r;
	vec3 shadedPosition = position;
	shadedPosition.y = center;

	gl_Position = mvp * vec4(shadedPosition, 1);

	// store linear depth
	depth = gl_Position.z / gl_Position.w;
}
