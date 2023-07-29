#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

layout (location = 0) in vec3 color;

layout (location = 0) out vec4 outColor;

float FogDensity = 0.04;

void main() 
{
	float depth = gl_FragCoord.z / gl_FragCoord.w - 0.2;
  
    // Calculate the fog factor
    float fogFactor = exp(-depth * FogDensity);
    fogFactor = clamp(fogFactor, 0.0, 1.0);


	outColor = vec4(color, fogFactor);
}
