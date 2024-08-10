#version 450

#include "lib/common.glsl"

IN(0, vec2, texcoord)
OUT(0, vec4, color)
SAMPLER(0, diffuse)

void main() 
{
	out_color = vec4(texture(tex_diffuse, in_texcoord).rgb, 1);
}
