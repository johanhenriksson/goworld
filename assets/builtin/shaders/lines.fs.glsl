#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec3, color)
OUT(0, vec4, color)

float FogDensity = 0.04;

void main() 
{
	float depth = gl_FragCoord.z / gl_FragCoord.w - 0.2;
  
    // Calculate the fog factor
    float fogFactor = exp(-depth * FogDensity);
    fogFactor = clamp(fogFactor, 0.0, 1.0);

	out_color = vec4(in_color, fogFactor);
}
