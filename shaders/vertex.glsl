#version 150 core

in vec2 vertex;
uniform mat4 mvmatrix;
uniform mat4 pmatrix;

void main() {
	gl_Position = pmatrix * mvmatrix * vec4(vertex, 0.0, 1.0);
}
