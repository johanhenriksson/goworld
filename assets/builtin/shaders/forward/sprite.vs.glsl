#version 450

#include "lib/common.glsl"
#include "lib/forward_vertex.glsl"

CAMERA(0, camera)
OBJECT(1, object)

// Attributes
IN(0, vec3, position)
IN(1, float, tex_x)
IN(2, vec3, normal)
IN(3, float, tex_y)
IN(4, vec4, color)

void main() 
{
	out_object = gl_InstanceIndex;

	vec3 center = (object.model * vec4(0, 0, 0, 1.0)).xyz;
	vec3 lookDirection = normalize(center - camera.Eye.xyz);
	vec3 up = vec3(0, 1, 0);

	vec3 right = normalize(cross(up, lookDirection));

	// if Y is not locked, calculate up vector
	// up = normalize(cross(lookDirection, right));

	out_world_position = center + 
		right * in_position.x +
		up * in_position.y;

	// texture coords
	out_color.xy = vec2(in_tex_x, in_tex_y);

	// gbuffer view position
	out_view_position = (camera.View * vec4(out_world_position.xyz, 1.0)).xyz;

	// world normal is always facing the camera
	out_world_normal = normalize(camera.Eye.xyz - center);

	// vertex clip space position
	gl_Position = camera.Proj * vec4(out_view_position, 1);
}
