#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"
#include "lib/forward_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)


VERTEX_BUFFER(Vertex)
INDEX_BUFFER(uint)

void main() 
{
	out_object = get_object_index();

	Vertex v = get_vertex_indexed(object.vertexBuffer, object.indexBuffer);

	// texture coords
	out_color.xy = vec2(v.tex_x, v.tex_y);

	// gbuffer view position
	out_view_position = (camera.View * object.model * vec4(v.position, 1)).xyz;
	out_world_position = (object.model * vec4(v.position, 1)).xyz;

	// world normal
	out_world_normal = normalize((object.model * vec4(v.normal, 0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
