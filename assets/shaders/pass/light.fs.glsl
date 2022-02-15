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
uniform mat4 cameraInverse; // inverse view projection matrix
uniform mat4 viewInverse;     // inverse view matrix
uniform mat4 light_vp;     // world to light space

uniform Light light;     // uniform light data
uniform float shadow_strength;
uniform float shadow_bias;
uniform bool soft_shadows = true;

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
    if (distanceToLight > light.Range) {
        return 0.0;
    }

    /* calculate normal coefficient */
    float normalCoef = max(0.0, dot(normal, surfaceToLight));

    /* light attenuation as a function of range and distance */
    float attenuation = light.attenuation.Constant +
                        light.attenuation.Linear * distanceToLight +
                        light.attenuation.Quadratic * pow(distanceToLight, 2);
    attenuation = 1.0 / attenuation;

    /* multiply and return light contribution */
    return normalCoef * attenuation;
}

/* samples the shadow map at the given world space coordinates */
float sampleShadowmap(sampler2D shadowmap, vec3 position, float bias) {
    /* world -> light clip coords */
    vec4 light_clip_pos = light_vp * vec4(position, 1);

    /* convert light clip to light ndc by dividing by W, then map values to 0-1 */
    vec3 light_ndc_pos = (light_clip_pos.xyz / light_clip_pos.w) * 0.5 + 0.5;

    /* depth of position in light space */
    float z = light_ndc_pos.z;
    if (z > 1) {
        return 0;
    }

    // todo: implement angle-dependent bias
    // float bias = max(0.05 * (1.0 - dot(normal, lightDir)), 0.005);  

    float shadow = 0.0;
    if (soft_shadows) {
        vec2 texelSize = 1.0 / textureSize(shadowmap, 0);
        for(int x = -1; x <= 1; ++x) {
            for(int y = -1; y <= 1; ++y) {
                float pcf_depth = texture(shadowmap, light_ndc_pos.xy + vec2(x, y) * texelSize).r; 
                shadow += z + bias > pcf_depth ? 1.0 : 0.0;        
            }    
        }
        shadow /= 9.0;
    }
    else {
        /* sample shadow map depth */
        float depth = texture(shadowmap, light_ndc_pos.xy).r;
        if (depth < (z + shadow_bias)) {
            shadow = 1.0; 
        }
    }

    return shadow * shadow_strength;
}

void main() {
    float depth = texture(tex_depth, texcoord0).r;

    // avoids lighting the backdrop.
    // perform this check early to avoid unnecessary work
    if (depth == 0.0) {
        discard;
    }

    // unpack data from geometry buffer
    vec4 t = texture(tex_diffuse, texcoord0);
    vec3 diffuseColor = t.rgb;
    float occlusion = t.a;

    // sample normal vector and transform it into world space
    vec3 normalEncoded = texture(tex_normal, texcoord0).xyz; // normals [0,1]
    vec3 viewNormal = normalize(2.0 * normalEncoded - 1); // normals [-1,1] 
    vec4 worldNormal = viewInverse * vec4(viewNormal, 0);
    vec3 normal = normalize(worldNormal.xyz);

    // calculate world position from depth map and the inverse camera view projection
    // why do we do this when we have a position buffer? :/
    vec3 position = positionFromDepth(depth);

    // now we should be able to calculate the position in light space!
    // since the directional light matrix is orthographic the z value will
    // correspond to depth, so we can test Z against the shadow map depth buffer
    // from the shadow pass! woop

    // calculate contribution from the light source
    float contrib = 0.0;
    float shadow = 1.0;
    if (light.Type == AMBIENT_LIGHT) {
        contrib = 1.0;
    }
    else if (light.Type == DIRECTIONAL_LIGHT) {
        // directional lights store the direction in the position uniform
        // i.e. the light coming from the position, shining towards the origin
        vec3 surfaceToLight = normalize(light.Position);
        contrib = max(dot(surfaceToLight, normal), 0.0);

        //float bias = max(0.05 * (1.0 - contrib), 0.005);  
        float bias = contrib * shadow_bias;

        // experimental shadows
        shadow = sampleShadowmap(tex_shadow, position, bias);
    }
    else if (light.Type == POINT_LIGHT) {
        // calculate light vector & distance
        vec3 surfaceToLight = light.Position - position;
        float distanceToLight = length(surfaceToLight);
        surfaceToLight = normalize(surfaceToLight);
        contrib = calculatePointLightContrib(surfaceToLight, distanceToLight, normal);
    }

    // calculate light color
    vec3 lightColor = light.Color * light.Intensity * contrib * shadow * occlusion;

    // mix with diffuse
    lightColor *= diffuseColor;

    // lightColor *= mix(1, ssao, ssao_amount);

    // write fragment color & restore depth buffer
    color = vec4(lightColor,  1.0);

    gl_FragDepth = depth;
}
