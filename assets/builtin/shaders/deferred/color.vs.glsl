#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

CAMERA(0, camera)
OBJECT(1, object, gl_InstanceIndex)

// Attributes
IN(0, vec3, position)
IN(1, vec3, normal)
IN(2, vec2, tex)
IN(3, vec4, color)

// Varyings
OUT(0, flat uint, object)
OUT(1, vec3, normal)
OUT(2, vec3, position)
OUT(3, vec4, color)

out gl_PerVertex {
	vec4 gl_Position;   
};

void main() 
{
	out_object = gl_InstanceIndex;
	mat4 mv = camera.View * object.model;

	// gbuffer diffuse
	out_color = vec4(in_color, 1);

	// gbuffer position
	out_position = (mv * vec4(in_position.xyz, 1.0)).xyz;

	// gbuffer view space normal
	out_normal = normalize((mv * vec4(in_normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
