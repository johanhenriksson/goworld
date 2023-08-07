struct Quad {
	vec2 min; // top left
	vec2 max; // bottom right
	vec2 uv_min; // top left uv
	vec2 uv_max; // bottom right uv
	vec4 color[4];
	float zindex;
	float corner_radius;
	float edge_softness;
	float border;
	uint texture;
};

UNIFORM(0, config, {
	vec2 resolution;
	float zmax;
})
STORAGE_BUFFER(1, Quad, quads)

float RoundedRectSDF(vec2 sample_pos, vec2 rect_center, vec2 rect_half_size, float r) {
	vec2 d2 = (abs(rect_center - sample_pos) - rect_half_size + vec2(r, r));
	return min(max(d2.x, d2.y), 0.0) + length(max(d2, 0.0)) - r;
}
