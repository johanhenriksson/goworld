
#define AMBIENT_LIGHT 0
#define POINT_LIGHT 1
#define DIRECTIONAL_LIGHT 2

const int SHADOW_CASCADES = 4;

layout (binding = 6) uniform sampler2D[] shadowmaps;

// these should be parameters
float shadow_power = 60;
float shadow_bias = 0.005;
float sample_radius = 1;
float normal_offset = 0.1;
int shadow_samples = 1;

struct Attenuation {
	float Constant;
	float Linear;
	float Quadratic;
};

struct Light {
	mat4 ViewProj[SHADOW_CASCADES];
	int Shadowmap[SHADOW_CASCADES];
	float Distance[SHADOW_CASCADES];

	vec4 Color;
	vec4 Position;
	int Type;
	float Intensity;
	float Range;
	Attenuation Attenuation;
};


// transforms ndc -> depth texture space
const mat4 biasMat = mat4( 
	0.5, 0.0, 0.0, 0.0,
	0.0, 0.5, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.0, 1.0 
);

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

float blendCascades(Light light, vec3 position, float depth, float bias, float blendRange, int numSamples, float sampleRadius) {
    // determine the cascade index
    int cascadeIndex = 0;
    for (int i = 0; i < SHADOW_CASCADES; ++i) {
        if (depth < light.Distance[i]) {
            cascadeIndex = i;
            break;
        }
    }
    sampleRadius *= cascadeIndex + 1;

    float shadowCurrent = sampleShadowmapPCF(shadowmaps[light.Shadowmap[cascadeIndex]], light.ViewProj[cascadeIndex], position, shadow_bias, numSamples, sampleRadius);

    // sample previous cascade
    if (cascadeIndex > 0) {
        // 4. Blend shadow results
        float cascadeStart = light.Distance[cascadeIndex - 1];
        float cascadeEnd = light.Distance[cascadeIndex];
        float blendFactor = smoothstep(cascadeStart, cascadeStart + blendRange, depth);

        if (blendFactor > 0) {
			float shadowPrev = sampleShadowmapPCF(shadowmaps[light.Shadowmap[cascadeIndex - 1]], light.ViewProj[cascadeIndex - 1], position, shadow_bias, numSamples-1, sampleRadius);
			return mix(shadowPrev, shadowCurrent, blendFactor);
        }
    }

    return shadowCurrent;
}

