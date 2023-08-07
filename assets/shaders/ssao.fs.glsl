#version 450 core
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

IN(0, vec2, texcoord)
OUT(0, float, ssao)

#define KERNEL_SIZE 32
layout (std140, binding = 0) uniform Params {
    mat4 Projection;
    vec4 Kernel[KERNEL_SIZE];
    int Samples;
    float Scale;
    float Radius;
    float Bias;
    float Power;
};

SAMPLER(1, position)
SAMPLER(2, normal)
SAMPLER(3, noise)

void main()
{
    vec2 noiseSize = vec2(textureSize(tex_noise, 0));
    vec2 outputSize = vec2(textureSize(tex_position, 0)) / Scale;
    vec2 noiseScale = outputSize / noiseSize;

    // get input vectors from gbuffer & noise texture
    vec3 fragPos = texture(tex_position, in_texcoord).xyz;
    vec3 normalEncoded = texture(tex_normal, in_texcoord).xyz;
    vec3 normal = unpack_normal(normalEncoded);

    // discard gbuffer entries without normal data
    if (normalEncoded == vec3(0)) {
        out_ssao = 1;
        return;
    }

    vec3 randomVec = texture(tex_noise, in_texcoord * noiseScale).xyz;

    // create TBN change-of-basis matrix: from tangent-space to view-space
    vec3 tangent = normalize(randomVec - normal * dot(randomVec, normal));
    vec3 bitangent = cross(normal, tangent);
    mat3 TBN = mat3(tangent, bitangent, normal);

    // iterate over the sample kernel and calculate occlusion factor
    float occlusion = 0.0;
    for(int i = 0; i < Samples; ++i)
    {
        // get sample position
        vec3 sampleVec = TBN * Kernel[i].xyz; // from tangent to view-space
        sampleVec = fragPos + sampleVec * Radius; 
        
        // project sample position (to sample texture) (to get position on screen/texture)
        vec4 offset = vec4(sampleVec, 1.0);
        offset = Projection * offset; // from view to clip-space
        offset.xyz /= offset.w; // perspective divide, clip -> NDC
        offset.xyz = offset.xyz * 0.5 + 0.5; // transform to range 0.0 - 1.0
        
        // get sample depth - i.e. the Z component of the sampled position in view space
        float sampleDepth = texture(tex_position, offset.xy).z;

        // range check & accumulate
        float rangeCheck = smoothstep(0.0, 1.0, Radius / abs(fragPos.z - sampleDepth));
        occlusion += (sampleDepth <= sampleVec.z - Bias ? 1.0 : 0.0) * rangeCheck;
    }
    occlusion = 1.0 - (occlusion / Samples);
    occlusion = pow(occlusion, Power);
    
    out_ssao = occlusion;
}
