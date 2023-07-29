// Standard Uniforms used in Deferred & Forward passes

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

layout (binding = 2) uniform sampler2D[] Textures;
