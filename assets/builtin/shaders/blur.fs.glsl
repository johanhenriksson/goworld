#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec2, texcoord)
OUT(0, vec4, color)
SAMPLER(0, input)

void main()
{             
    vec2 texelSize = 1.0 / vec2(textureSize(tex_input, 0));
    float result = 0.0;
    for (int x = -2; x < 2; ++x) 
    {
        for (int y = -2; y < 2; ++y) 
        {
            vec2 offset = vec2(float(x), float(y)) * texelSize;
            result += texture(tex_input, in_texcoord + offset).r;
        }
    }
    result = result / (4.0 * 4.0);
    out_color = vec4(result, result, result, 1);
}
