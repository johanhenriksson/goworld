#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec2, texcoord)
OUT(0, vec4, color)
SAMPLER(0, input)
SAMPLER(1, lut)

#define MAXCOLOR 15.0 
#define COLORS 16.0
#define WIDTH 256.0 
#define HEIGHT 16.0

vec3 lookup_color(sampler2D lut, vec3 clr) {
    float cell = clr.b * MAXCOLOR;

    float cell_l = floor(cell); 
    float cell_h = ceil(cell);
    
    float half_px_x = 0.5 / WIDTH;
    float half_px_y = 0.5 / HEIGHT;
    float r_offset = half_px_x + clr.r / COLORS * (MAXCOLOR / COLORS);
    float g_offset = half_px_y + clr.g * (MAXCOLOR / COLORS);
    
    vec2 lut_pos_l = vec2(cell_l / COLORS + r_offset, 1 - g_offset); 
    vec2 lut_pos_h = vec2(cell_h / COLORS + r_offset, 1 - g_offset);

    vec3 graded_color_l = texture(lut, lut_pos_l).rgb; 
    vec3 graded_color_h = texture(lut, lut_pos_h).rgb;

    return mix(graded_color_l, graded_color_h, fract(cell));
}

void main() {
    // todo: expose as uniform setting
    float exposure = 1.0;

    // get input color
    vec3 hdrColor = texture(tex_input, in_texcoord).rgb;

    // exposure tone mapping
    vec3 mapped = vec3(1.0) - exp(-hdrColor * exposure);

    // gamma correction
    vec3 corrected = pow(mapped, vec3(1/gamma));

    // color grading
    vec3 graded = lookup_color(tex_lut, corrected);

    // return
    out_color = vec4(graded, 1);
}
