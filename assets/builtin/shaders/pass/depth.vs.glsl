#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

OUT(0, vec3, position)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

CAMERA(0, camera)
OBJECT(1, object, gl_InstanceIndex)

VERTEX_BUFFER(Vertex)
INDEX_BUFFER(uint)

void main() 
{
	// load vertex data
	Vertex v = get_vertex_indexed(object.vertexPtr, object.indexPtr);

	mat4 mv = camera.View * object.model;

	// gbuffer view position
	// todo: can this be removed? probably just a waste of bandwidth
	out_position = (mv * vec4(v.position, 1.0)).xyz;

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
