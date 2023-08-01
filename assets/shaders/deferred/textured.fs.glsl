#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/fragment.glsl"

void main() 
{
	vec2 texcoord0 = color0.xy;

	uint texture0 = objects.item[objectIndex].textures[0];
	diffuse = vec4(texture(Textures[texture0], texcoord0).rgb, 1);

	normal = pack_normal(normal0);

	position = vec4(position0, 1);
}
