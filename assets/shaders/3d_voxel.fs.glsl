#version 330

uniform float ambient;
uniform float lightIntensity;
uniform sampler2D tex0;

in vec2 texcoord;
in vec3 worldNormal;
in vec3 L;
in vec3 V;
in float lightDistance;

layout(location=1) out vec4 outputColor;
layout(location=0) out vec4 outputSpecular;
layout(location=2) out vec4 outputNormal;

void main() {
    vec4 diffuse = texture(tex0, texcoord);
    float i = lightIntensity / pow(lightDistance, 2);
    float light = max(0, dot(L, worldNormal)) * i;

    light = max(ambient, light);
    vec4 color = diffuse * light;
    color.w = 1.0;

    outputColor = color;
    outputNormal = vec4((worldNormal + 1) / 2, 1);
    outputSpecular = vec4(0,1,0,1);
}
