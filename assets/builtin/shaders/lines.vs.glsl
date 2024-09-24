#version 450

#include "lib/common.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
Object object = objects.item[gl_InstanceIndex];

IN(0, vec3, position)
IN(4, vec4, color)
OUT(0, vec3, color)

out gl_PerVertex 
{
    vec4 gl_Position;   
};

void main() 
{
    out_color = in_color.rgb;

	mat4 mvp = camera.ViewProj * object.model;
	gl_Position = mvp * vec4(in_position, 1);
}
