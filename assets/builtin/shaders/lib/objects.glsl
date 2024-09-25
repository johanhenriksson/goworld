// Standard engine vertex format
// Size: 48 bytes
struct Vertex {
	vec3 position;
	vec3 normal;
	vec2 tex;
	vec4 color;
};

struct Object {
	mat4 model;
	uint textures[MAX_TEXTURES];

	uint64_t vertexPtr;
	uint64_t indexPtr;
};

#define get_object_index() (gl_InstanceIndex)

#define get_vertex_indexed(vertexPtr, indexPtr) (VertexBuffer(vertexPtr)[IndexBuffer(indexPtr)[gl_VertexIndex].index].vertex)
#define get_vertex(vertexPtr) (VertexBuffer(vertexPtr)[gl_VertexIndex].vertex)

#define VERTEX_BUFFER(VertexType) layout(buffer_reference, scalar, buffer_reference_align=4) readonly buffer VertexBuffer { VertexType vertex; };
#define INDEX_BUFFER(IndexType) layout(buffer_reference, scalar, buffer_reference_align=4) readonly buffer IndexBuffer { IndexType index; };
