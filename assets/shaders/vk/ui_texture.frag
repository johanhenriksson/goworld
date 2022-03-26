#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

layout (binding = 0) uniform sampler2D[] textures;

layout (push_constant) uniform constants
{
    mat4 Viewport;
    mat4 Model;
    int Texture;
} element;

layout (location = 0) in vec2 out_uv;
layout (location = 1) in vec4 out_color;

layout (location = 0) out vec4 frag_color;

void main() {
    vec2 uv = out_uv;

    // uv's outside 0.0-1.0 should be transparent/discarded
    if (all(lessThan(uv, vec2(0.0))) && all(greaterThan(uv, vec2(1.0)))) {
        discard;
    }

    vec4 tint = vec4(1);
    if (out_color.a > 0) {
        tint = out_color;
    }

    vec4 texcolor = texture(textures[element.Texture], uv);
    frag_color = texcolor.rgba * tint;
}
