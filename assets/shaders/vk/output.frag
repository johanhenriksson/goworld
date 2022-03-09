#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

layout (set = 0, binding = 2) uniform sampler2D diffuse;

layout (location = 0) in vec2 texcoord;

// Return Output
layout (location = 0) out vec4 outFragColor;

void main() 
{
    vec3 color = texture(diffuse, texcoord).rgb;
    outFragColor = vec4(color, 1.0);
}
