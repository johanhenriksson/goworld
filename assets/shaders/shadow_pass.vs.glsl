#version 330

#define DIRECTIONAL_LIGHT 2

struct Light {
    Attenuation attenuation;
    vec3 Color;
    vec3 Position;
    float Range;
    int Type;
};

uniform Light light;
uniform mat4 projection; // orthographic for directional lights, otherwise its a perspective matrix

void main() {

    if (light.Type == DIRECTIONAL_LIGHT) {

    }

}
