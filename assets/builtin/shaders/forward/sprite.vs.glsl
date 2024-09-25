#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

OUT(0, flat uint, object)
OUT(1, vec4, color)
OUT(2, vec2, texcoord)
OUT(3, vec3, view_position)
OUT(4, vec3, world_normal)
OUT(5, vec3, world_position)

out gl_PerVertex {
	vec4 gl_Position;   
};

CAMERA(0, camera)
OBJECT(1, object, gl_InstanceIndex)

VERTEX_BUFFER(Vertex)
INDEX_BUFFER(uint)

void main() 
{
	out_object = get_object_index();

	// load vertex data
	Vertex v = get_vertex_indexed(object.vertexPtr, object.indexPtr);

	vec3 center = (object.model * vec4(0, 0, 0, 1.0)).xyz;
	vec3 lookDirection = normalize(center - camera.Eye.xyz);
	vec3 up = vec3(0, 1, 0);

	vec3 right = normalize(cross(up, lookDirection));

	// if Y is not locked, calculate up vector
	// up = normalize(cross(lookDirection, right));

	out_world_position = center + 
		right * v.position.x +
		up * v.position.y;

	// texture & color
	out_texcoord = v.tex;
	out_color = v.color;

	// gbuffer view position
	out_view_position = (camera.View * vec4(out_world_position.xyz, 1.0)).xyz;

	// world normal is always facing the camera
	out_world_normal = normalize(camera.Eye.xyz - center);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
