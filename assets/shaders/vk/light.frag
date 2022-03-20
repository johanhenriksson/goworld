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
} camera;

layout (binding = 5) uniform sampler2D[] shadowmaps;

layout(push_constant) uniform constants
{
    mat4 ViewProj;
    vec4 Color;
    vec4 Position;
    int Type;
    int Shadowmap;
    float Range;
    float Intensity;
    Attenuation Attenuation;
} light;

layout (input_attachment_index = 0, binding = 0) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 1) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 2) uniform subpassInput tex_position;
layout (input_attachment_index = 3, binding = 3) uniform subpassInput tex_depth;

layout (location = 0) out vec4 color;

float shadow_bias = 0.005;
bool soft_shadows = true;
float shadow_strength = 0.75;

vec3 getWorldPosition() {
    /* world position */
    vec4 pos_ws = camera.ViewInv * vec4(subpassLoad(tex_position).xyz, 1);
    return pos_ws.xyz / pos_ws.w;
}

vec3 getWorldNormal() {
    // sample normal vector and transform it into world space
    vec3 viewNormal = normalize(2.0 * subpassLoad(tex_normal).rgb - 1); // normals [-1,1] 
    vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
    return normalize(worldNormal.xyz);
}

float sampleShadowmap(sampler2D shadowmap, vec3 position, float bias) {
    /* world -> light clip coords */
    vec4 light_pos = light.ViewProj * vec4(position, 1);
    light_pos = light_pos / light_pos.w;

    /* convert light clip to light ndc by dividing by W, then map values to 0-1 */
    light_pos.st = light_pos.st * 0.5 + 0.5;

    /* depth of position in light space */
    float z = light_pos.z;
    if (z < -1 || z > 1) {
        return 0.0;
    }

    float shadow = 0.0;
    if (soft_shadows) {
        vec2 texelSize = 1.0 / textureSize(shadowmap, 0);
        for(int x = -1; x <= 1; ++x) {
            for(int y = -1; y <= 1; ++y) {
                float pcf_depth = texture(shadowmap, light_pos.st + vec2(x, y) * texelSize).r; 
                shadow += pcf_depth > (z - bias) ? 1.0 : 0.0;        
            }    
        }
        shadow /= 9.0;
    }
    else {
        /* sample shadow map depth */
        float depth = texture(shadowmap, light_pos.st).r;
        if (depth > (z - bias)) {
            shadow = 1.0; 
        }
    }

    return shadow * shadow_strength;
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
    // discard empty fragments
    float depth = subpassLoad(tex_depth).r;
    if (depth == 1) {
        discard;
    }

    // unpack data from geometry buffer
    vec4 t = subpassLoad(tex_diffuse);
    vec3 diffuseColor = t.rgb;
    float occlusion = t.a;

    vec3 position = getWorldPosition();
    vec3 normal = getWorldNormal();

    // now we should be able to calculate the position in light space!
    // since the directional light matrix is orthographic the z value will
    // correspond to depth, so we can test Z against the shadow map depth buffer
    // from the shadow pass! woop

    // calculate contribution from the light source
    float contrib = 0.0;
    float shadow = 1.0;
    if (light.Type == AMBIENT_LIGHT) {
        contrib = 1;
    }
    else if (light.Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        // i.e. the light coming from the position, shining towards the origin
        vec3 surfaceToLight = normalize(light.Position.xyz);
        contrib = max(dot(surfaceToLight, normal), 0.0);

        // angle-dependent bias
        // shadow_bias = max(shadow_bias * (1.0 - dot(normal, surfaceToLight)), 0.0005);  

        // experimental shadows
        if (light.Shadowmap > 0) {
            shadow = sampleShadowmap(shadowmaps[light.Shadowmap], position, shadow_bias);
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
    color = vec4(lightColor,  1.0);
}
