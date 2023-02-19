#version 330
  
in vec2 texcoord0;
layout(location=0) out vec4 color;
  
uniform sampler2D tex_input;

void main()
{             
    vec2 texelSize = 1.0 / vec2(textureSize(tex_input, 0));
    float result = 0.0;
    for (int x = -2; x < 2; ++x) 
    {
        for (int y = -2; y < 2; ++y) 
        {
            vec2 offset = vec2(float(x), float(y)) * texelSize;
            result += texture(tex_input, texcoord0 + offset).r;
        }
    }
    result = result / (4.0 * 4.0);
    color = vec4(result, result, result, 1);
}