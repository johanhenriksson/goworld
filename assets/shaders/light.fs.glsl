#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

#define AMBIENT_LIGHT 0
#define POINT_LIGHT 1
#define DIRECTIONAL_LIGHT 2

struct Attenuation {
	float Constant;
	float Linear;
	float Quadratic;
};

layout (std140, binding = 4) uniform Camera {
	mat4 Proj;
	mat4 View;
	mat4 ViewProj;
	mat4 ProjInv;
	mat4 ViewInv;
	mat4 ViewProjInv;
	vec3 Eye;
	vec3 Forward;
} camera;

// could be a pipeline parameter
const int cascades = 4;

struct Light {
	mat4 ViewProj[cascades];
	int Shadowmap[cascades];
	float Distance[cascades];
};

layout (std430, binding = 5) readonly buffer LightBuffer {
	Light lights[];
} ssbo;

layout (binding = 6) uniform sampler2D[] shadowmaps;

layout(push_constant) uniform constants
{
	mat4 ViewProj;
	vec4 Color;
	vec4 Position;
	int Type;
	int Index;
	float Range;
	float Intensity;
	Attenuation Attenuation;
} light;

layout (input_attachment_index = 0, binding = 0) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 1) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 2) uniform subpassInput tex_position;
layout (input_attachment_index = 3, binding = 3) uniform subpassInput tex_depth;

layout (location = 0) out vec4 color;

const mat4 biasMat = mat4( 
	0.5, 0.0, 0.0, 0.0,
	0.0, 0.5, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.0, 1.0 
);

float shadow_bias = 0.005;

vec3 getWorldPosition(vec3 viewPos) {
	// transform view space to world space
	vec4 pos_ws = camera.ViewInv * vec4(viewPos, 1);
	return pos_ws.xyz / pos_ws.w;
}

float getDepth() {
	return subpassLoad(tex_position).z;
}

vec3 getWorldNormal() {
	// sample normal vector and transform it into world space
	vec3 viewNormal = normalize(2.0 * subpassLoad(tex_normal).rgb - 1); // normals [-1,1] 
	vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
	return normalize(worldNormal.xyz);
}

float sampleShadowmap(sampler2D shadowmap, mat4 viewProj, vec3 position) {
	vec4 shadowCoord = biasMat * viewProj * vec4(position, 1);

	float shadow = 1.0;
	if (shadowCoord.z > -1.0 && shadowCoord.z < 1.0) {
		float dist = texture(shadowmap, shadowCoord.st).r;
		if (shadowCoord.w > 0 && dist < shadowCoord.z - shadow_bias) {
			shadow = 0;
		}
	}
	return shadow;
}

/* calculates lighting contribution from a point light source */
float calculatePointLightContrib(vec3 surfaceToLight, float distanceToLight, vec3 normal) {
	if (distanceToLight > light.Range) {
		return 0.0;
	}

	/* calculate normal coefficient */
	float normalCoef = max(0.0, dot(normal, surfaceToLight));

	/* light attenuation as a function of range and distance */
	float attenuation = light.Attenuation.Constant +
						light.Attenuation.Linear * distanceToLight +
						light.Attenuation.Quadratic * pow(distanceToLight, 2);
	attenuation = 1.0 / attenuation;

	/* multiply and return light contribution */
	return normalCoef * attenuation;
}

void main() {
	vec3 viewPos = subpassLoad(tex_position).xyz;

	// unpack data from geometry buffer
	vec4 t = subpassLoad(tex_diffuse);
	vec3 diffuseColor = t.rgb;
	float occlusion = t.a;

	vec3 position = getWorldPosition(viewPos);
	vec3 normal = getWorldNormal();

	// calculate contribution from the light source
	float contrib = 0.0;
	float shadow = 1.0;
	if (light.Type == AMBIENT_LIGHT) {
		contrib = 1;
	}
	else if (light.Type == DIRECTIONAL_LIGHT) {
		// directional lights store the direction in the position uniform
		// i.e. the light coming from the position, shining towards the origin
		vec3 lightDir = normalize(light.Position.xyz);
		vec3 surfaceToLight = -lightDir;
		contrib = max(dot(surfaceToLight, normal), 0.0);

		// experimental shadows
		if (light.Index > 0) {
			// find light struct
			Light dirlight = ssbo.lights[light.Index];

			// pick cascade index
			int index = 0;
			for(int i = 0; i < cascades; i++) {
				if (viewPos.z < dirlight.Distance[i]) {
					index = i;
					break;
				}
			}

			shadow = sampleShadowmap(shadowmaps[dirlight.Shadowmap[index]], dirlight.ViewProj[index], position);
		}
	}
	else if (light.Type == POINT_LIGHT) {
		// calculate light vector & distance
		vec3 surfaceToLight = light.Position.xyz - position;
		float distanceToLight = length(surfaceToLight);
		surfaceToLight = normalize(surfaceToLight);
		contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
	} 

	vec3 lightColor = light.Color.rgb * light.Intensity * contrib * shadow * occlusion;
	lightColor *= diffuseColor;

	// lightColor *= mix(1, ssao, ssao_amount);

	// write fragment color & restore depth buffer
	color = vec4(lightColor, 1.0);
}
