#version 450 core

// Input: Texture coords
layout (location = 0) in vec2 texcoord0;

// Return Output
layout (location = 0) out vec4 out_ssao;

layout (std140, binding = 0) uniform Params {
    mat4 Proj;
    vec3 Kernel[32];
} params;

layout (binding = 1) uniform sampler2D tex_position; // position gbuffer
layout (binding = 2) uniform sampler2D tex_normal; // normal gbuffer
layout (binding = 3) uniform sampler2D tex_noise; // noise texture

// parameters (you'd probably want to use them as uniforms to more easily tweak the effect)
int kernel_size = 32;
float radius = 0.4;
float bias = 0.025;
float power = 1.2;
int scale = 1;


void main()
{
    vec2 noiseScale = vec2(textureSize(tex_position, 0)) / 4.0 / scale;

    // get input vectors from gbuffer & noise texture
    vec3 fragPos = texture(tex_position, texcoord0).xyz;
    vec3 normalEncoded = texture(tex_normal, texcoord0).xyz; // normals [0,1]
    vec3 normal = normalize(2.0 * normalEncoded - 1); // normals [-1,1] 

    vec3 randomVec = texture(tex_noise, texcoord0 * noiseScale).xyz;

    // create TBN change-of-basis matrix: from tangent-space to view-space
    vec3 tangent = normalize(randomVec - normal * dot(randomVec, normal));
    vec3 bitangent = cross(normal, tangent);
    mat3 TBN = mat3(tangent, bitangent, normal);

    // iterate over the sample kernel and calculate occlusion factor
    float occlusion = 0.0;
    for(int i = 0; i < kernel_size; ++i)
    {
        // get sample position
        vec3 sampleVec = TBN * params.Kernel[i]; // from tangent to view-space
        sampleVec = fragPos + sampleVec * radius; 
        
        // project sample position (to sample texture) (to get position on screen/texture)
        vec4 offset = vec4(sampleVec, 1.0);
        offset = params.Proj * offset; // from view to clip-space
        offset.xyz /= offset.w; // perspective divide, clip -> NDC
        offset.xyz = offset.xyz * 0.5 + 0.5; // transform to range 0.0 - 1.0
        
        // get sample depth - i.e. the Z component of the sampled position in view space
        float sampleDepth = texture(tex_position, offset.xy).z;

        // range check & accumulate
        float rangeCheck = smoothstep(0.0, 1.0, radius / abs(fragPos.z - sampleDepth));
        if (sampleDepth <= sampleVec.z - bias) {
            occlusion += 1.0 * rangeCheck;
        }
    }
    occlusion = 1.0 - (occlusion / kernel_size);
    occlusion = pow(occlusion, power);
    
    out_ssao = vec4(vec3(occlusion), 1);
}
