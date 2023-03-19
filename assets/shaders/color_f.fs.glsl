#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

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

// Varying
layout (location = 0) in vec4 color0;
layout (location = 1) in vec3 normal0;
layout (location = 2) in vec3 position0;
layout (location = 3) in vec3 wnormal;

// Return Output
layout (location = 0) out vec4 diffuse;
layout (location = 1) out vec4 normal;
layout (location = 2) out vec4 position;

void main() 
{
    vec3 lightDir = normalize(camera.Forward);
    vec3 surfaceToLight = -lightDir;
    float contrib = max(dot(surfaceToLight, wnormal), 0.2);

    diffuse = vec4(color0.rgb * contrib, color0.a);

    vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
    normal = pack_normal;

    position = vec4(position0, 1);
}
