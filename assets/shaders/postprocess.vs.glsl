#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec3, position)
IN(1, vec2, texcoord)
OUT(0, vec2, texcoord)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	out_texcoord = in_texcoord;
	gl_Position = vec4(in_position, 1);
}
