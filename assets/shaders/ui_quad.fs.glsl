#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Varying
layout (location = 0) in vec4 color0;

// Return Output
layout (location = 0) out vec4 output_color;

layout (binding = 2) uniform sampler2D[] Textures;

void main() 
{
    output_color = color0;
}
