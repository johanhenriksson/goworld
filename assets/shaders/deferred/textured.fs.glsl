#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable
#extension GL_GOOGLE_include_directive : enable

#include "../lib/uniforms.glsl"
#include "../lib/fragment.glsl"

void main() 
{
	vec2 texcoord0 = color0.xy;

	uint texture0 = ssbo.objects[objectIndex].textures[0];
	diffuse = vec4(texture(Textures[texture0], texcoord0).rgb, 1);

	vec4 pack_normal = vec4((normalize(normal0) + 1.0) / 2.0, 1);
	normal = pack_normal;

	position = vec4(position0, 1);
}
