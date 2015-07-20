#version 330

uniform sampler2D tex0;

in vec2 texcoord;
in vec3 frag_normal;

out vec4 outputColor;

void main() {
    outputColor = 0.1 * vec4(frag_normal,1) + texture(tex0, texcoord); 
}
