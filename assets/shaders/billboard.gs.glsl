#version 330

layout (points) in;
layout (triangle_strip) out;
layout (max_vertices = 4) out;

uniform mat4 m;
uniform mat4 vp;
uniform vec3 eye;

out vec2 uv;

void main()
{
    vec3 pos = gl_in[0].gl_Position.xyz;

    vec3 toCamera = normalize(eye - pos);
    vec3 up = vec3(0.0, 1.0, 0.0);
    vec3 right = cross(toCamera, up);

    pos -= (right * 0.5);
    gl_Position = vp * vec4(pos, 1.0);
    uv = vec2(0.0, 0.0);
    EmitVertex();

    pos.y += 1.0;
    gl_Position = vp * vec4(pos, 1.0);
    uv = vec2(0.0, 1.0);
    EmitVertex();

    pos.y -= 1.0;
    pos += right;
    gl_Position = vp * vec4(pos, 1.0);
    uv = vec2(1.0, 0.0);
    EmitVertex();

    pos.y += 1.0;
    gl_Position = vp * vec4(pos, 1.0);
    uv = vec2(1.0, 1.0);
    EmitVertex();

    EndPrimitive();
} 