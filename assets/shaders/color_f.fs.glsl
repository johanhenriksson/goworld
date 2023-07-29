#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_GOOGLE_include_directive : enable

#include "lib/uniforms.glsl"
#include "lib/fragment.glsl"

// Varying
layout (location = 4) in vec3 wnormal;

float gamma = 2.2;

void main() 
{
    vec3 lightDir = normalize(camera.Forward);
    vec3 surfaceToLight = -lightDir;
    float contrib = max(dot(surfaceToLight, wnormal), 0.2);

    // gamma correct & write fragment
	vec3 linearColor = pow(color0.rgb, vec3(gamma));
    diffuse = vec4(linearColor * contrib, color0.a);

    // update gbuffer
    vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
    normal = pack_normal;

    position = vec4(position0, 1);
}
