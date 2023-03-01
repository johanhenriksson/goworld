#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

// Varying
layout (location = 0) in vec4 color0;
layout (location = 1) in vec2 uv0;
layout (location = 2) flat in uint texture0;
layout (location = 3) in vec2 pos;
layout (location = 4) flat in vec2 center;
layout (location = 5) flat in vec2 half_size;
layout (location = 6) flat in float corner_radius;
layout (location = 7) flat in float edge_softness;

// Return Output
layout (location = 0) out vec4 output_color;

layout (binding = 2) uniform sampler2D[] Textures;

float RoundedRectSDF(vec2 sample_pos, vec2 rect_center, vec2 rect_half_size, float r)
{
  vec2 d2 = (abs(rect_center - sample_pos) -
             rect_half_size +
             vec2(r, r));
  return min(max(d2.x, d2.y), 0.0) + length(max(d2, 0.0)) - r;
}

void main() 
{
    // shrink the rectangle's half-size that is used for distance calculations 
    // otherwise the underlying primitive will cut off the falloff too early.
    vec2 softness_padding = vec2(max(0, edge_softness*2-1),
                                 max(0, edge_softness*2-1));

    // sample distance to rect at position
    float dist = RoundedRectSDF(pos,
                                center,
                                half_size-softness_padding,
                                corner_radius);

    // map distance to a blend factor
    float sdf_factor = 1 - smoothstep(0, 2*edge_softness, dist);

    vec4 sample0 = texture(Textures[texture0], uv0);
    output_color = color0 * sample0 * sdf_factor;
}
