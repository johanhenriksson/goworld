#version 330 core

layout(location=0) out vec4 out_ssao;

in vec2 texcoord0;

uniform mat4 projection;

uniform sampler2D tex_position; // position gbuffer
uniform sampler2D tex_normal; // normal gbuffer
uniform sampler2D tex_noise; // noise texture

uniform vec3 samples[32];

// parameters (you'd probably want to use them as uniforms to more easily tweak the effect)
uniform int kernel_size = 32;
uniform float radius = 0.4;
uniform float bias = 0.025;
uniform float power = 1.2;
uniform int scale = 1;


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
        vec3 sample = TBN * samples[i]; // from tangent to view-space
        sample = fragPos + sample * radius; 
        
        // project sample position (to sample texture) (to get position on screen/texture)
        vec4 offset = vec4(sample, 1.0);
        offset = projection * offset; // from view to clip-space
        offset.xyz /= offset.w; // perspective divide, clip -> NDC
        offset.xyz = offset.xyz * 0.5 + 0.5; // transform to range 0.0 - 1.0
        
        // get sample depth - i.e. the Z component of the sampled position in view space
        float sampleDepth = texture(tex_position, offset.xy).z;

        // range check & accumulate
        float rangeCheck = smoothstep(0.0, 1.0, radius / abs(fragPos.z - sampleDepth));
        if (sampleDepth <= sample.z - bias) {
            occlusion += 1.0 * rangeCheck;
        }
    }
    occlusion = 1.0 - (occlusion / kernel_size);
    occlusion = pow(occlusion, power);
    
    out_ssao = vec4(vec3(occlusion), 1);
}