#version 330

uniform float ambient;
uniform float lightIntensity;
uniform sampler2D tex0;

in vec2 texcoord;
in vec3 worldNormal;
in vec3 L;
in vec3 V;
in float lightDistance;

out vec4 outputColor;

void main() {
    vec4 diffuse = texture(tex0, texcoord);
    float i = lightIntensity / pow(lightDistance, 2);
    float light = max(0, dot(L, worldNormal)) * i;

    light = max(ambient, light);
    vec4 color = diffuse * light;
    color.w = 1.0;

    outputColor = color;
}
