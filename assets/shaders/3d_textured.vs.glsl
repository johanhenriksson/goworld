#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vertex;
in vec2 texCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = texCoord;
    gl_Position = projection * camera * model * vec4(vertex, 1);
}
