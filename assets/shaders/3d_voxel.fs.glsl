#version 330

uniform sampler2D tex0;

in vec2 texcoord;
in vec3 worldNormal;
in vec3 L;
in vec3 V;
in float lightDistance;

out vec4 outputColor;

void main() {
    float ambient = 0.25;
    float intensity = 15.0;

    vec4 diffuse = texture(tex0, texcoord);
    float i = intensity / pow(lightDistance, 2);
    float light = max(0, dot(L, worldNormal)) * i;

    light = max(ambient, light);
    vec4 color = diffuse * light;
    color.w = 1.0;

    outputColor = color;
}
