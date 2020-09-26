#version 330

uniform sampler2D sprite;

in vec2 uv;
out vec4 color;

void main()
{
    color = texture(sprite, uv) + vec4(1,0,0,1);
} 