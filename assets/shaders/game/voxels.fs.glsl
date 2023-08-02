#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/material.glsl"
#include "lib/fragment.glsl"

void main() 
{
	diffuse = vec4(color0.rgb * color0.a, 1);
	normal = pack_normal(normal0);
	position = vec4(position0, 1);
}
