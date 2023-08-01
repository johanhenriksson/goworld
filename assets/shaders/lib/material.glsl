// Standard Uniforms used in Deferred & Forward passes

struct ObjectData{
	mat4 model;
	uint textures[4];
};

layout (binding = 1) readonly buffer ObjectBuffer {
	ObjectData item[];
} objects;

layout (binding = 2) readonly buffer LightBuffer {
	LightSettings settings;
	float[LIGHT_PADDING] _padding;
	Light item[];
} lights;

layout (binding = 3) uniform sampler2D[] Textures;
