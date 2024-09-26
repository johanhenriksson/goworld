#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

OUT(0, vec3, color)

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

    out_color = v.color.rgb;

	mat4 mvp = camera.ViewProj * object.model;
	gl_Position = mvp * vec4(v.position, 1);
}
