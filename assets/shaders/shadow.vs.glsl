#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)

// Attributes
IN(0, vec3, position)
OUT(0, float, depth)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

void main() 
{
	mat4 mvp = camera.ViewProj * objects.item[gl_InstanceIndex].model;
	gl_Position = mvp * vec4(in_position, 1);

	// store linear depth
	out_depth = gl_Position.z / gl_Position.w;
}
