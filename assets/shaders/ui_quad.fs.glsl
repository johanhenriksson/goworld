#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

// Varying
layout (location = 0) in vec4 color0;
layout (location = 1) flat in uint texture0;
layout (location = 2) in vec2 uv0;

// Return Output
layout (location = 0) out vec4 output_color;

layout (binding = 2) uniform sampler2D[] Textures;

void main() 
{
    vec4 sample0 = texture(Textures[texture0], uv0);
    output_color = color0 * sample0;
}
