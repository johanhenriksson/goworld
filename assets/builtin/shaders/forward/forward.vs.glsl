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

	// texture coords
	out_texcoord = v.tex;
	out_color = v.color;

	// gbuffer view position
	out_view_position = (camera.View * object.model * vec4(v.position, 1)).xyz;
	out_world_position = (object.model * vec4(v.position, 1)).xyz;

	// world normal
	out_world_normal = normalize((object.model * vec4(v.normal, 0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
