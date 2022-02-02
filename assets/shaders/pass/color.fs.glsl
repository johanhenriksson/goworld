#version 330

in vec2 texcoord0;
layout(location=0) out vec4 color;

uniform sampler2D tex_input; // source image
uniform sampler2D tex_lut; // color lookup table
uniform sampler2D tex_ssao;
uniform float gamma = 2.2;

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
    // get input color
    vec3 src = texture(tex_input, texcoord0.xy).rgb;


    // color grading
    vec3 graded = lookup_color(tex_lut, src);

    float ssao = texture(tex_ssao, texcoord0.xy).r;
    graded *= ssao;

    // gamma correction
    vec3 corrected = pow(graded, vec3(1/gamma));

    // return
    color = vec4(corrected, 1);
}