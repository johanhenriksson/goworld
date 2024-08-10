#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec3, position)
OUT(0, vec4, position)

void main() 
{
    out_position = vec4(in_position, 1);
}
