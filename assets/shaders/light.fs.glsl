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

// transforms ndc -> depth texture space
const mat4 biasMat = mat4( 
	0.5, 0.0, 0.0, 0.0,
	0.0, 0.5, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.0, 1.0 
);

float shadow_power = 60;
float shadow_bias = 0.005;
float sample_radius = 1;
float normal_offset = 0.1;
int shadow_samples = 2;
const bool debug = false;

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

float sampleShadowmap(sampler2D shadowmap, mat4 viewProj, vec3 position, float bias) {
	vec4 shadowCoord = biasMat * viewProj * vec4(position, 1);

	float shadow = 1.0;
	if (shadowCoord.z > -1.0 && shadowCoord.z < 1.0 && shadowCoord.w > 0) {
		float dist = texture(shadowmap, shadowCoord.st).r;
		float actual = exp(shadow_power * shadowCoord.z - bias) / exp(shadow_power);

		if (dist < actual) {
			shadow = 0;
		}
	}
	return shadow;
}

float sampleShadowmapPCF(sampler2D shadowmap, mat4 viewProj, vec3 position, float bias, int numSamples, float sampleRadius) {
	if (numSamples <= 0) {
		return sampleShadowmap(shadowmap, viewProj, position, bias);
	}

	vec4 shadowCoord = biasMat * viewProj * vec4(position, 1);
    shadowCoord = shadowCoord / shadowCoord.w;

    float shadow = 1.0;
    if (shadowCoord.z > -1.0 && shadowCoord.z < 1.0 && shadowCoord.w > 0) {
        vec2 shadowMapSize = textureSize(shadowmap, 0).xy;
        vec2 texelSize = 1.0 / shadowMapSize;
        float actual = exp(shadow_power * (shadowCoord.z - bias)) / exp(shadow_power);

        float count = 0.0;
        for (int x = -numSamples; x <= numSamples; ++x) {
            for (int y = -numSamples; y <= numSamples; ++y) {
                vec2 offset = vec2(float(x), float(y)) * texelSize * sampleRadius;
                float dist = texture(shadowmap, shadowCoord.st + offset).r;

                // Compare the difference between exponential depth values
                if (dist - actual < 0) {
                    count += 1.0;
                }
            }
        }

        shadow = 1 - count / float((2 * numSamples + 1) * (2 * numSamples + 1));
    }
    return shadow;
}

float blendCascades(int shadowIndices[cascades], mat4 viewProj[cascades], float cascadeSplits[cascades], vec3 position, float depth, float bias, float blendRange, int numSamples, float sampleRadius) {
    // determine the cascade index
    int cascadeIndex = 0;
    for (int i = 0; i < cascades; ++i) {
        if (depth < cascadeSplits[i]) {
            cascadeIndex = i;
            break;
        }
    }
    sampleRadius *= cascadeIndex + 1;

    float shadowCurrent = sampleShadowmapPCF(shadowmaps[shadowIndices[cascadeIndex]], viewProj[cascadeIndex], position, shadow_bias, numSamples, sampleRadius);

    // sample previous cascade
    if (cascadeIndex > 0) {
        // 4. Blend shadow results
        float cascadeStart = cascadeSplits[cascadeIndex - 1];
        float cascadeEnd = cascadeSplits[cascadeIndex];
		float blendDistance = cascadeSplits[cascadeIndex-1] * blendRange;
        float blendFactor = smoothstep(cascadeStart, cascadeStart + blendDistance, depth);

        if (blendFactor > 0) {
			float shadowPrev = sampleShadowmapPCF(shadowmaps[shadowIndices[cascadeIndex - 1]], viewProj[cascadeIndex - 1], position, shadow_bias, numSamples-1, sampleRadius);
			return mix(shadowPrev, shadowCurrent, blendFactor);
        }
    }

    return shadowCurrent;
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

			float bias = shadow_bias * max(0.0, 1.0 - dot(normal, lightDir));
			position += normal * normal_offset;
			shadow = blendCascades(dirlight.Shadowmap, dirlight.ViewProj, dirlight.Distance, position, viewPos.z, bias, 0.3, shadow_samples, sample_radius);

			if (debug) {
				int index = -1;
				for(int i = 0; i < cascades; i++) {
					if (viewPos.z < dirlight.Distance[i]) {
						index = i;
						break;
					}
				}
				diffuseColor = mix(vec3(0,1,0), vec3(1,0,0), float(index) / (cascades - 1));
			}
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
