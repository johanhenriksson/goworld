#version 330

uniform sampler2D tex_diffuse; // diffuse
uniform sampler2D tex_normal; // normal
uniform sampler2D tex_depth; // depth

uniform mat4 cameraInverse;

// light data
uniform vec3 l_position;
uniform vec3 l_intensity;
uniform float l_attenuation_const;
uniform float l_attenuation_linear;
uniform float l_attenuation_quadratic;
uniform float l_range;

in vec2 texcoord0;

out vec4 color;

vec3 positionFromDepth(float depth) {
    /* homogenous coords */
    float xhs = 2 * texcoord0.x - 1;
    float yhs = 2 * texcoord0.y - 1;
    float zhs = 2 * depth - 1;

    /* homogenous vector */
    vec4 pos_hs = vec4(xhs, yhs, zhs, 1) / gl_FragCoord.w;

    /* world position */
    vec4 pos_ws = cameraInverse * pos_hs;
    return pos_ws.xyz / pos_ws.w;
}

vec4 calculatePointLight(vec3 surfaceToLight, float distanceToLight, vec3 normal) {
    float diffuseCoefficient = max(0.0, dot(normal, surfaceToLight));
    float attenuation = l_attenuation_const +
                        l_attenuation_linear * distanceToLight +
                        l_attenuation_quadratic * distanceToLight * distanceToLight;
    attenuation = 1 / attenuation;
    attenuation *= clamp(pow(1.0 - pow(distanceToLight / l_range, 4), 2), 0, 1);

    vec4 diffuse = vec4(0.0);
    diffuse.rgb = l_intensity * diffuseCoefficient * attenuation;
    diffuse.a = 1.0;

    return diffuse;
}

void main() {
    vec3 diffuseColor = texture(tex_diffuse, texcoord0).rgb;
    vec3 normalEncoded = texture(tex_normal, texcoord0).xyz;
    float depth = texture(tex_depth, texcoord0).r;

    vec3 position = positionFromDepth(depth);
    vec3 normal = normalize(2.0 * normalEncoded - 1);

    vec3 surfaceToLight = normalize(l_position - position);
    float distanceToLight = length(l_position - position);
    
    vec4 lightColor = calculatePointLight(surfaceToLight, distanceToLight, normal);

    vec3 gamma = vec3(1.0 / 2.2);

    vec4 phat = vec4(diffuseColor + normalEncoded, depth);

    color = (0.1 + 
            lightColor) * vec4(pow(diffuseColor, gamma),1)
            //vec4(distanceToLight * 0.1 * position, 1)
            + 0.001 * phat;
}
