#version 330

in vec3 fragColor;

out vec4 outputColor;

void main() {
    outputColor = vec4(fragColor,1);
}
