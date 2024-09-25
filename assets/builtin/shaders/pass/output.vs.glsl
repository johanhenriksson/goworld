#version 450

#include "lib/common.glsl"

IN(0, vec3, position)
IN(2, vec2, tex)
OUT(0, vec2, texcoord)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	out_texcoord = in_tex;
	gl_Position = vec4(in_position, 1);
}
