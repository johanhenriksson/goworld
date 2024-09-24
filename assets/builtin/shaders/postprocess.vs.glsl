#version 450

#include "lib/common.glsl"

IN(0, vec3, position)
IN(1, float, tex_x)
IN(3, float, tex_y)
OUT(0, vec2, texcoord)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	out_texcoord = vec2(in_tex_x, in_tex_y);
	gl_Position = vec4(in_position, 1);
}
