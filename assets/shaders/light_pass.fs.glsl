#version 330

#define AMBIENT_LIGHT 0
#define POINT_LIGHT 1
#define DIRECTIONAL_LIGHT 2

struct Attenuation {
    float Constant;
    float Linear;
    float Quadratic;
};

struct Light {
    Attenuation attenuation;
    vec3 Color;
    vec3 Position;
    float Range;
    float Intensity;
    int Type;
};

uniform sampler2D tex_diffuse;  // diffuse
uniform sampler2D tex_shadow;  // shadow map
uniform sampler2D tex_normal; // normal
uniform sampler2D tex_depth; // depth
uniform sampler2D tex_occlusion; // ssao
uniform mat4 cameraInverse; // inverse view projection matrix
uniform mat4 light_vp;     // world to light space
uniform mat4 viewInverse;     // projection matrix

uniform Light light;     // uniform light data
uniform float shadow_strength;
uniform float shadow_bias;
uniform float ssao_amount = 1.0;

in vec2 texcoord0;

out vec4 color;

/* dark mathemagic - Translates from clip space back into world space */
vec3 positionFromDepth(float depth) {
    /* clip coords */
    float xhs = 2 * texcoord0.x - 1;
    float yhs = 2 * texcoord0.y - 1;
    float zhs = 2 * depth - 1;

    /* homogenous clip vector */
    vec4 pos_hs = vec4(xhs, yhs, zhs, 1) / gl_FragCoord.w;

    /* world position */
    vec4 pos_ws = cameraInverse * pos_hs;
    return pos_ws.xyz / pos_ws.w;
}

/* calculates lighting contribution from a point light source */
float calculatePointLightContrib(vec3 surfaceToLight, float distanceToLight, vec3 normal) {
    /* calculate normal coefficient */
    float normalCoef = max(0.0, dot(normal, surfaceToLight));

    /* light attenuation as a function of range and distance */
    float attenuation = light.attenuation.Constant +
                        light.attenuation.Linear * distanceToLight +
                        light.attenuation.Quadratic * pow(distanceToLight, 2);
    attenuation = light.Range / attenuation;

    /* multiply and return light contribution */
    return normalCoef * attenuation;
}

/* samples the shadow map at the given world space coordinates */
float sampleShadowmap(sampler2D shadowmap, vec3 position) {
    /* world -> light clip coords */
    vec4 light_clip_pos = light_vp * vec4(position, 1);

    /* convert light clip to light ndc by dividing by W, then map values to 0-1 */
    vec3 light_ndc_pos = (light_clip_pos.xyz / light_clip_pos.w) * 0.5 + 0.5;

    /* depth of position in light space */
    float z = light_ndc_pos.z;

    /* sample shadow map depth */
    float depth = texture(shadowmap, light_ndc_pos.xy).r;

    /* shadow test */
    if (depth < (z - shadow_bias)) {
        return 1.0 - shadow_strength;
    }

    return 1.0;
}

void main() {
    /* unpack data from geometry buffer */
    vec4 t = texture(tex_diffuse, texcoord0);
    vec3 diffuseColor = t.rgb;
    float occlusion = t.a;

    // sample normal vector and transform it into world space
    vec3 normalEncoded = texture(tex_normal, texcoord0).xyz; // normals [0,1]
    vec3 viewNormal = normalize(2.0 * normalEncoded - 1); // normals [-1,1] 
    vec4 worldNormal = viewInverse * vec4(viewNormal, 0);
    vec3 normal = normalize(worldNormal.xyz);

    /* calculate world position from depth map and the inverse camera view projection */
    float depth = texture(tex_depth, texcoord0).r;
    vec3 position = positionFromDepth(depth);

    float ssao = texture(tex_occlusion, texcoord0).r;

    // now we should be able to calculate the position in light space!
    // since the directional light matrix is orthographic the z value will
    // correspond to depth, so we can test Z against the shadow map depth buffer
    // from the shadow pass! woop

    /* calculate contribution from the light source */
    float contrib = 0.0;
    float shadow = 1.0;
    if (light.Type == AMBIENT_LIGHT) {
        contrib = 1.0;
    }
    else if (light.Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        vec3 dir = normalize(light.Position);
        contrib = max(dot(dir, normal), 0.0);

        // experimental shadows
        shadow = sampleShadowmap(tex_shadow, position);
    }
    else if (light.Type == POINT_LIGHT) {
        /* calculate light vector & distance */
        vec3 surfaceToLight = light.Position - position;
        float distanceToLight = length(surfaceToLight);
        surfaceToLight = normalize(surfaceToLight);
        contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
    }

    // avoids lighting the backdrop.
    // probably inefficient though, consider another solution.
    if (depth == 1.0) {
        contrib = 0;
    }

    /* calculate light color */
    vec3 lightColor = light.Color * light.Intensity * contrib * shadow * occlusion;

    /* mix with diffuse */
    lightColor *= diffuseColor;

    lightColor *= mix(1, ssao, ssao_amount);

    /* write fragment color & restore depth buffer */
    color = vec4(lightColor,  1.0);

    gl_FragDepth = depth;
}
