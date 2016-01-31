#version 330

uniform sampler2D tex0;

in vec2 texcoord;
in vec3 worldNormal;

out vec3 out_diffuse;
layout(location=1) out vec3 out_normal;

void main() {
    out_diffuse = texture(tex0, texcoord).rgb;
    out_normal = 0.001*(worldNormal+1)/2;
}
