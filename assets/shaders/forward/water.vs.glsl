#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/forward_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)

// Attributes
IN(0, vec3, position)
IN(1, vec3, normal)
IN(2, vec2, texcoord)

const float epsilon = 0.05;
const float waveHeight = 0.1;
const float waveFrequency = 2;
const float waveSpeed = 1.1;

float getWaveHeight(vec3 position) {
    float wave1 = sin(position.x * waveFrequency + camera.Time * waveSpeed) * 0.5;
    float wave2 = sin(position.z * waveFrequency * 0.8 + camera.Time * waveSpeed * 1.2) * 0.4;
    float wave3 = sin((position.x + position.z) * waveFrequency * 1.2 + camera.Time * waveSpeed * 0.8) * 0.3;
    return (wave1 + wave2 + wave3) * waveHeight;
}

vec3 calculateWaveNormal(vec3 position, float height, float epsilon) {
    // Calculate heights at neighboring points
    vec3 tangentX = vec3(epsilon, getWaveHeight(position + vec3(epsilon, 0, 0)) - height, 0);
    vec3 tangentZ = vec3(0, getWaveHeight(position + vec3(0, 0, epsilon)) - height, epsilon);
    
    // Calculate normal using cross product
    return normalize(cross(tangentZ, tangentX));
}

void main() 
{
	out_object = gl_InstanceIndex;

	// texture coords
	out_color.xy = in_texcoord;

	// gbuffer view position
	out_world_position = (object.model * vec4(in_position.xyz, 1.0)).xyz;

	// add wave height
    float height = getWaveHeight(out_world_position);
	out_world_position.y += height;

	out_view_position = (camera.View * vec4(out_world_position, 1)).xyz;

	// world normal
	vec3 waveNormal = calculateWaveNormal(out_world_position, height, epsilon);
	out_world_normal = normalize((object.model * vec4(waveNormal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
