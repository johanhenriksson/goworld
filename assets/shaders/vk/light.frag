#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

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

layout (std140, binding = 5) uniform Light {
    mat4 Proj;
    mat4 View;
    mat4 ViewProj;
    vec4 Color;
    vec4 Position;
    int Type;
    float Range;
    float Intensity;
    int Shadows;
    Attenuation Attenuation;
} lights[10];

layout(push_constant) uniform constants
{
	int lightId;
} push;

layout (input_attachment_index = 0, binding = 0) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 1) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 2) uniform subpassInput tex_position;
layout (input_attachment_index = 3, binding = 3) uniform subpassInput tex_depth;

layout (location = 0) out vec4 color;

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

/* calculates lighting contribution from a point light source */
float calculatePointLightContrib(vec3 surfaceToLight, float distanceToLight, vec3 normal) {
    if (distanceToLight > lights[push.lightId].Range) {
        return 0.0;
    }

    /* calculate normal coefficient */
    float normalCoef = max(0.0, dot(normal, surfaceToLight));

    /* light attenuation as a function of range and distance */
    float attenuation = lights[push.lightId].Attenuation.Constant +
                        lights[push.lightId].Attenuation.Linear * distanceToLight +
                        lights[push.lightId].Attenuation.Quadratic * pow(distanceToLight, 2);
    attenuation = 1.0 / attenuation;

    /* multiply and return light contribution */
    return normalCoef * attenuation;
}

void main() {
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
    if (lights[push.lightId].Type == AMBIENT_LIGHT) {
        contrib = 1.0;
    }
    else if (lights[push.lightId].Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        // i.e. the light coming from the position, shining towards the origin
        vec3 surfaceToLight = normalize(lights[push.lightId].Position.xyz);
        contrib = max(dot(surfaceToLight, normal), 0.0);

        // angle-dependent bias
        // float bias = max(shadow_bias * (1.0 - dot(normal, surfaceToLight)), 0.0005);  

        // experimental shadows
        // shadow = sampleShadowmap(tex_shadow, position, bias);
    }
    else if (lights[push.lightId].Type == POINT_LIGHT) {
        // calculate light vector & distance
        vec3 surfaceToLight = lights[push.lightId].Position.xyz - position;
        float distanceToLight = length(surfaceToLight);
        surfaceToLight = normalize(surfaceToLight);
        contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
    } 

    vec3 lightColor = lights[push.lightId].Color.rgb * lights[push.lightId].Intensity * contrib * shadow * occlusion;
    lightColor *= diffuseColor;

    // lightColor *= mix(1, ssao, ssao_amount);

    // write fragment color & restore depth buffer
    color = vec4(lightColor,  1.0);
}
