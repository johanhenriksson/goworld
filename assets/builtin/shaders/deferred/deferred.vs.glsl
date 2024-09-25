#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

// Varyings
OUT(0, flat uint, object)
OUT(1, vec3, position)
OUT(2, vec3, normal)
OUT(3, vec4, color)
OUT(4, vec2, texcoord)

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
	Vertex vtx = get_vertex_indexed(object.vertexPtr, object.indexPtr);

	mat4 mv = camera.View * object.model;

	// textures
	out_texcoord = vtx.tex;
	out_color = vtx.color;

	// gbuffer position
	out_position = (mv * vec4(vtx.position, 1)).xyz;

	// gbuffer normal
	out_normal = normalize((mv * vec4(vtx.normal, 0.0)).xyz);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_position, 1);
}
