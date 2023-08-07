#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)

// Attributes
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

// Varyings
layout (location = 0) out vec3 normal0;
layout (location = 1) out vec3 position0;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	mat4 m = objects.item[gl_InstanceIndex].model;
	mat4 mv = camera.View * m;

	// gbuffer view position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
