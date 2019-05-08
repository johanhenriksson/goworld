#version 330
  
in vec2 texcoord0;
layout(location=0) out vec4 color;
  
uniform sampler2D tex_input;

const int size = 4;

const float weight[5] = float[] (0.227027, 0.1945946, 0.1216216, 0.054054, 0.016216);

void main()
{             
    vec2 tex_offset = 1.0 / textureSize(tex_input, 0); // gets size of single texel
    vec3 result = texture(tex_input, texcoord0).rgb * weight[0]; // current fragment's contribution
    for(int i = 1; i < 5; ++i)
    {
        result += texture(tex_input, texcoord0 + vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
        result += texture(tex_input, texcoord0 - vec2(tex_offset.x * i, 0.0)).rgb * weight[i];
    }

    color = vec4(result, 1);
}