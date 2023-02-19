#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

layout (push_constant) uniform constants
{
    mat4 Viewport;
    mat4 Model;
    int Texture;
} element;

layout (location = 0) in vec3 position;
layout (location = 1) in vec4 color_0;
layout (location = 2) in vec2 texcoord_0;

layout (location = 0) out vec2 out_uv;
layout (location = 1) out vec4 out_color;

void main() {
    out_uv      = texcoord_0;
    out_color   = color_0;
    gl_Position = element.Viewport * element.Model * vec4(position, 1);
}
