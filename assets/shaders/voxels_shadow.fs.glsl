#version 450

layout (location = 0) in float depth;

layout (location = 0) out vec4 fragColor;

float shadow_power = 60;

void main() 
{
	fragColor = vec4(0);

	// exponential depth
	gl_FragDepth = exp(shadow_power * depth) / exp(shadow_power);
}
