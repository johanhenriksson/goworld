#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_vertex.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)

// Attributes
IN(0, vec3, position)
IN(1, vec3, normal)
IN(2, vec2, texcoord)
IN(3, vec4, weights)

OUT(4, vec4, weights)

void main() 
{
	out_object = gl_InstanceIndex;
	mat4 mv = camera.View * objects.item[out_object].model;

	// textures
	out_color.xy = in_texcoord;

	// pass texture weights
	out_weights = in_weights;

	// gbuffer position
	out_position = (mv * vec4(in_position, 1)).xyz;

	// gbuffer normal
	out_normal = normalize((mv * vec4(in_normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}