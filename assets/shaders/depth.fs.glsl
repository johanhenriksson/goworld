#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

layout (location = 0) in vec3 normal0;
layout (location = 1) in vec3 position0;

// Return Output
layout (location = 0) out vec4 normal;
layout (location = 1) out vec4 position;

void main() 
{
    normal = pack_normal(normal0);
    position = vec4(position0, 1);
}
