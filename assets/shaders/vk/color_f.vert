#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Attributes
layout (location = 0) in vec3 inPos;
layout (location = 1) in vec3 inColor;

// layout (push_constant) uniform Push {
// } push;

layout (std140, set = 0, binding = 0) uniform UBO {
	mat4 proj;
	mat4 view;
} ubo;

struct ObjectData{
	mat4 model;
};

//all object matrices
layout(std140, set = 1, binding = 0) readonly buffer ObjectBuffer{
	ObjectData objects[];
} ssbo;

// Varyings
layout (location = 0) out vec3 outColor;

out gl_PerVertex 
{
    vec4 gl_Position;   
};


void main() 
{
	outColor = inColor;
	gl_Position = ubo.proj * ubo.view * ssbo.objects[gl_InstanceIndex].model * vec4(inPos.xyz, 1.0);
}