#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/fragment.glsl"

// Varying
layout (location = 4) in vec3 wnormal;

void main() 
{
    vec3 lightDir = normalize(camera.Forward);
    vec3 surfaceToLight = -lightDir;
    float contrib = max(dot(surfaceToLight, wnormal), 0.2);

    // gamma correct & write fragment
	vec3 linearColor = pow(color0.rgb, vec3(gamma));
    diffuse = vec4(linearColor * contrib, color0.a);

    // update gbuffer
    normal = pack_normal(normal0);
    position = vec4(position0, 1);
}
