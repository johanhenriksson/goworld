#define POINT_LIGHT 1
#define DIRECTIONAL_LIGHT 2

#define SHADOW_CASCADES 4

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
	uint Type;
	float Intensity;
	float Range;
	Attenuation Attenuation;
};

struct LightSettings {
	vec4 AmbientColor;
	float AmbientIntensity;
	int Count;
	int ShadowSamples;
	float ShadowSampleRadius;
	float ShadowBias;
	float NormalOffset;
};

#define LIGHT_PADDING 76
#define LIGHT_BUFFER(idx,name) layout (binding = idx) readonly buffer LightBuffer { LightSettings settings; float[LIGHT_PADDING] _padding; Light item[]; } name;

const float SHADOW_POWER = 60;

// transforms ndc -> depth texture space
const mat4 biasMat = mat4( 
	0.5, 0.0, 0.0, 0.0,
	0.0, 0.5, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.0, 1.0 
);

//
// Lighting functions
//

float _shadow_texture(uint index, vec2 point);
vec2 _shadow_size(uint index);
float sampleShadowmap(uint shadowmap, mat4 viewProj, vec3 position, float bias);
float sampleShadowmapPCF(uint shadowmap, mat4 viewProj, vec3 position, LightSettings settings);
float blendCascades(Light light, vec3 position, float depth, float blendRange, LightSettings settings);
float calculatePointLightContrib(Light light, vec3 surfaceToLight, float distanceToLight, vec3 normal);
vec3 ambientLight(LightSettings settings, float occlusion);
vec3 calculateLightColor(Light light, vec3 position, vec3 normal, float depth, LightSettings settings);

float sampleShadowmap(uint shadowmap, mat4 viewProj, vec3 position, float bias) {
	vec4 shadowCoord = biasMat * viewProj * vec4(position, 1);

	float shadow = 1.0;
	if (shadowCoord.z > -1.0 && shadowCoord.z < 1.0 && shadowCoord.w > 0) {
		float dist = _shadow_texture(shadowmap, shadowCoord.st);
		float actual = exp(SHADOW_POWER * shadowCoord.z - bias) / exp(SHADOW_POWER);

		if (dist < actual) {
			shadow = 0;
		}
	}
	return shadow;
}

float sampleShadowmapPCF(uint shadowmap, mat4 viewProj, vec3 position, LightSettings settings) {
	if (settings.ShadowSamples <= 0) {
		return sampleShadowmap(shadowmap, viewProj, position, settings.ShadowBias);
	}

	vec4 shadowCoord = biasMat * viewProj * vec4(position, 1);
    shadowCoord = shadowCoord / shadowCoord.w;

    float shadow = 1.0;
    if (shadowCoord.z > -1.0 && shadowCoord.z < 1.0 && shadowCoord.w > 0) {
        vec2 texelSize = 1.0 / _shadow_size(shadowmap);
        float actual = exp(SHADOW_POWER * (shadowCoord.z - settings.ShadowBias)) / exp(SHADOW_POWER);

        float count = 0.0;
        int numSamples = settings.ShadowSamples;
        for (int x = -numSamples; x <= numSamples; ++x) {
            for (int y = -numSamples; y <= numSamples; ++y) {
                vec2 offset = vec2(float(x), float(y)) * texelSize * settings.ShadowSampleRadius;
                float dist = _shadow_texture(shadowmap, shadowCoord.st + offset);

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

float blendCascades(Light light, vec3 position, float depth, float blendRange, LightSettings settings) {
    // determine the cascade index
    int cascadeIndex = 0;
    for (int i = 0; i < SHADOW_CASCADES; ++i) {
        if (depth < light.Distance[i]) {
            cascadeIndex = i;
            break;
        }
    }

    float shadowCurrent = sampleShadowmapPCF(light.Shadowmap[cascadeIndex], light.ViewProj[cascadeIndex], position, settings);

    // blend with previous cascade to get a smooth transition
    if (cascadeIndex > 0 && blendRange > 0) {
        float cascadeStart = light.Distance[cascadeIndex - 1];
        float cascadeEnd = light.Distance[cascadeIndex];
        blendRange *= cascadeIndex;
        float blendFactor = smoothstep(cascadeStart, cascadeStart + blendRange, depth);

        if (blendFactor > 0) {
			float shadowPrev = sampleShadowmapPCF(light.Shadowmap[cascadeIndex - 1], light.ViewProj[cascadeIndex - 1], position, settings);
			return mix(shadowPrev, shadowCurrent, blendFactor);
        }
    }

    return shadowCurrent;
}

/* calculates lighting contribution from a point light source */
float calculatePointLightContrib(Light light, vec3 surfaceToLight, float distanceToLight, vec3 normal) {
	if (distanceToLight > light.Range) {
		return 0.0;
	}

	// calculate normal coefficient
	float normalCoef = max(0.0, dot(normal, surfaceToLight));

	// light attenuation as a function of range and distance
	float attenuation = light.Attenuation.Constant +
						light.Attenuation.Linear * distanceToLight +
						light.Attenuation.Quadratic * pow(distanceToLight, 2);
	attenuation = 1.0 / attenuation;

	// multiply and return light contribution
	return normalCoef * attenuation;
}

vec3 ambientLight(LightSettings settings, float occlusion) {
	return settings.AmbientColor.rgb * settings.AmbientIntensity * occlusion;
}

vec3 calculateLightColor(Light light, vec3 position, vec3 normal, float depth, LightSettings settings) {
	float contrib = 0.0;
	float shadow = 1.0;
	if (light.Type == DIRECTIONAL_LIGHT) {
		// directional lights store the direction in the position uniform
		// i.e. the light coming from the position, shining towards the origin
		vec3 lightDir = normalize(light.Position.xyz);
		vec3 surfaceToLight = -lightDir;
		contrib = max(dot(surfaceToLight, normal), 0.0);

		float bias = settings.ShadowBias * max(0.0, 1.0 - dot(normal, lightDir));
		position += normal * settings.NormalOffset;
		shadow = blendCascades(light, position, depth, light.Range, settings);

#if DEBUG_CASCADES
		int index = -1;
		for(int i = 0; i < SHADOW_CASCADES; i++) {
			if (depth < light.Distance[i]) {
				index = i;
				break;
			}
		}
		return contrib * shadow * mix(vec3(0,1,0), vec3(1,0,0), float(index) / (SHADOW_CASCADES - 1));
#endif
	}
	else if (light.Type == POINT_LIGHT) {
		// calculate light vector & distance
		vec3 surfaceToLight = light.Position.xyz - position;
		float distanceToLight = length(surfaceToLight);
		surfaceToLight = normalize(surfaceToLight);
		contrib = calculatePointLightContrib(light, surfaceToLight, distanceToLight, normal);
	} 

	return light.Color.rgb * light.Intensity * contrib * shadow;
}
