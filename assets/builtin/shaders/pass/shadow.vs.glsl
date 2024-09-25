#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
Object object = objects.item[gl_InstanceIndex];

// Attributes
IN(0, vec3, position)
OUT(0, float, depth)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	mat4 mvp = camera.ViewProj * object.model;
	gl_Position = mvp * vec4(in_position, 1);

	// store linear depth
	out_depth = gl_Position.z / gl_Position.w;
}
