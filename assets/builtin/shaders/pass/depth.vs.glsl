#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
Object object = objects.item[gl_InstanceIndex];

// Attributes
layout (location = 0) in vec3 position;

// Varyings
layout (location = 0) out vec3 position0;

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	mat4 m = object.model;
	mat4 mv = camera.View * m;

	// gbuffer view position
	position0 = (mv * vec4(position.xyz, 1.0)).xyz;

	// vertex clip space position
	gl_Position = camera.Proj * vec4(position0, 1);
}
