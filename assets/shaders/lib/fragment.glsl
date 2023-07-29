// Common fragment shader code
// Varying
layout (location = 0) in vec4 color0;
layout (location = 1) in vec3 normal0;
layout (location = 2) in vec3 position0;
layout (location = 3) in flat uint objectIndex;

// Return Output
layout (location = 0) out vec4 diffuse;
layout (location = 1) out vec4 normal;
layout (location = 2) out vec4 position;
