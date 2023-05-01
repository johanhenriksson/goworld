#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in uint normal_id;
layout (location = 2) in vec3 color_0;
layout (location = 3) in float occlusion_0;

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
layout (location = 0) out vec3 color0;
layout (location = 1) out vec3 normal0;
layout (location = 2) out vec3 position0;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

const vec3 normals[7] = vec3[7] (
	vec3(0,0,0),  // normal 0 - undefined
	vec3(1,0,0),  // x+
	vec3(-1,0,0), // x-
	vec3(0,1,0),  // y+
	vec3(0,-1,0), // y-
	vec3(0,0,1),  // z+
	vec3(0,0,-1)  // z-
);

void main() 
{
	mat4 mv = camera.View * ssbo.objects[gl_InstanceIndex].model;

	// gbuffer diffuse
	color0 = color_0 * (1 - occlusion_0);

	// gbuffer position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	vec3 normal = normals[normal_id];
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
