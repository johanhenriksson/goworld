#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

// Uniforms

layout (binding = 2) uniform sampler2D[] Textures;

// Varying
layout (location = 0) in vec3 normal0;
layout (location = 1) in vec3 position0;
layout (location = 2) in flat uint texture0;
layout (location = 3) in vec2 texcoord0;

// Return Output
layout (location = 0) out vec4 diffuse;
layout (location = 1) out vec4 normal;
layout (location = 2) out vec4 position;

void main() 
{
	vec4 color0 = vec4(texture(Textures[texture0], texcoord0).rgb, 1);
	diffuse = color0;

	vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
	normal = pack_normal;

	normal = vec4(0,1,0,1);

	position = vec4(position0, 1);
}
