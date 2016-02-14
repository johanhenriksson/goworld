#version 330

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
    int Type;
};

uniform sampler2D tex_diffuse;  // diffuse
uniform sampler2D tex_shadow;  // shadow map
uniform sampler2D tex_normal; // normal
uniform sampler2D tex_depth; // depth
uniform mat4 cameraInverse; // inverse view projection matrix
uniform mat4 light_vp;     // world to light space

uniform Light light;     // uniform light data

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

vec3 gammaCorrect(vec3 color) {
    const vec3 gamma = vec3(1.0 / 2.2);
    return pow(color, gamma);
}

/* samples the shadow map at the given world space coordinates */
float sampleShadowmap(sampler2D shadowmap, vec3 position) {
    /* world -> light clip coords */
    vec4 light_clip_pos = light_vp * vec4(position,1);

    /* convert light clip to light ndc by dividing by W, then map values to 0-1 */
    vec3 light_ndc_pos = (light_clip_pos.xyz / light_clip_pos.w) * 0.5 + 0.5;

    /* depth of position in light space */
    float z = light_ndc_pos.z;

    /* sample shadow map depth */
    float depth = texture(shadowmap, light_ndc_pos.xy).r;

    /* shadow test */
    if (depth < (z + 0.0001)) {
        return 0.5;
    }

    return 1.0;
}

void main() {
    /* unpack data from geometry buffer */
    vec4 t = texture(tex_diffuse, texcoord0);
    vec3 diffuseColor = t.rgb;
    float occlusion = t.a;

    vec3 normalEncoded = texture(tex_normal, texcoord0).xyz; // normals [0,1]
    vec3 normal = normalize(2.0 * normalEncoded - 1); // normals [-1,1]

    float depth = texture(tex_depth, texcoord0).r;

    /* calculate world position from depth map and the inverse camera view projection */
    vec3 position = positionFromDepth(depth);

    // now we should be able to calculate the position in light space!
    // since the directional light matrix is orthographic the z value will
    // correspond to depth, so we can test Z against the shadow map depth buffer
    // from the shadow pass! woop

    /* calculate contribution from the light source */
    float contrib = 0.0;
    if (light.Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        vec3 dir = normalize(-light.Position);
        contrib = max(dot(dir, normal), 0.0);

        // experimental shadows
        float shadow = sampleShadowmap(tex_shadow, position);
        occlusion *= shadow;

    }
    else if (light.Type == POINT_LIGHT) {
        /* calculate light vector & distance */
        vec3 surfaceToLight = light.Position - position;
        float distanceToLight = length(surfaceToLight);
        surfaceToLight = normalize(surfaceToLight);
        contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
    }

    /* calculate light color */
    vec3 lightColor = light.Color * occlusion * contrib;

    /* add ambient light */
    const vec3 ambientColor = vec3(0.95, 1.0, 0.91);
    lightColor += 0.1 * ambientColor;

    /* mix with diffuse */
    lightColor *= diffuseColor;

    /* write fragment color & restore depth buffer */
    color = vec4(lightColor, 1.0);
    gl_FragDepth = depth;
}
