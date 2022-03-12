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

layout (binding = 3) uniform Camera {
    mat4 Proj;
    mat4 View;
    mat4 ViewProj;
    mat4 ProjInv;
    mat4 ViewInv;
    mat4 ViewProjInv;
    vec3 Eye;
} camera;

layout (binding = 4) uniform Light {
    vec3 Position;
    vec3 Color;
    float Range;
    float Intensity;
    int Type;
    mat4 Proj;
    mat4 View;
    mat4 ViewProj;
    Attenuation attenuation;
} lights[10];

layout(push_constant) uniform constants
{
	int lightId;
} push;

layout (input_attachment_index = 0, binding = 0) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 1) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 2) uniform subpassInput tex_position;

layout (location = 0) out vec4 color;

vec3 getWorldPosition() {
    /* world position */
    vec4 pos_ws = camera.ViewInv * vec4(subpassLoad(tex_position).xyz, 1);
    return pos_ws.xyz / pos_ws.w;
}

void main() {
    vec3 position = getWorldPosition();
    int light = push.lightId;

    // unpack data from geometry buffer
    vec4 t = subpassLoad(tex_diffuse);
    vec3 diffuseColor = t.rgb;
    float occlusion = 1; // t.a;

    // sample normal vector and transform it into world space
    vec3 viewNormal = normalize(2.0 * subpassLoad(tex_normal).rgb - 1); // normals [-1,1] 
    vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
    vec3 normal = normalize(worldNormal.xyz);

    // now we should be able to calculate the position in light space!
    // since the directional light matrix is orthographic the z value will
    // correspond to depth, so we can test Z against the shadow map depth buffer
    // from the shadow pass! woop

    // calculate contribution from the light source
    float contrib = 0.0;
    float shadow = 1.0;
    if (lights[light].Type == AMBIENT_LIGHT) {
        contrib = 1.0;
    }
    else if (lights[light].Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        // i.e. the light coming from the position, shining towards the origin
        vec3 surfaceToLight = normalize(lights[light].Position);
        contrib = max(dot(surfaceToLight, normal), 0.0);

        // angle-dependent bias
        // float bias = max(shadow_bias * (1.0 - dot(normal, surfaceToLight)), 0.0005);  

        // experimental shadows
        // shadow = sampleShadowmap(tex_shadow, position, bias);
    }
    // else if (light.Type == POINT_LIGHT) {
    //     // calculate light vector & distance
    //     vec3 surfaceToLight = light.Position - position;
    //     float distanceToLight = length(surfaceToLight);
    //     surfaceToLight = normalize(surfaceToLight);
    //     contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
    // }

    // calculate light color
    vec3 lightColor = lights[light].Color * lights[light].Intensity * contrib * shadow * occlusion;

    // mix with diffuse
    lightColor *= diffuseColor;
    // lightColor = normal;

    // lightColor *= mix(1, ssao, ssao_amount);

    // write fragment color & restore depth buffer
    color = vec4(lightColor,  1.0);
}
