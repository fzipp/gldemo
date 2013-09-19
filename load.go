package main

import (
	"errors"
	"github.com/fzipp/geom"
	"github.com/go-gl/gl"
	"io/ioutil"
)

func LoadShader(filename string, shaderType gl.GLenum) (gl.Shader, error) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	shader := gl.CreateShader(shaderType)
	shader.Source(string(code))
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != 1 {
		return shader, errors.New(shader.GetInfoLog())
	}
	return shader, nil
}

func LoadShaders(vertFile, fragFile string) (gl.Shader, gl.Shader, error) {
	vs, err := LoadShader(vertFile, gl.VERTEX_SHADER)
	if err != nil {
		return vs, 0, err
	}
	fs, err := LoadShader(fragFile, gl.FRAGMENT_SHADER)
	if err != nil {
		return vs, fs, err
	}
	return vs, fs, nil
}

func LoadProgram(shaders ...gl.Shader) (gl.Program, error) {
	p := gl.CreateProgram()
	for _, shader := range shaders {
		p.AttachShader(shader)
	}
	p.Link()
	if p.Get(gl.LINK_STATUS) != 1 {
		return p, errors.New(p.GetInfoLog())
	}
	p.Validate()
	if p.Get(gl.VALIDATE_STATUS) != 1 {
		return p, errors.New(p.GetInfoLog())
	}
	return p, nil
}

func LoadMatrix(m *geom.Mat4, p gl.Program, name string) {
	loc := p.GetUniformLocation(name)
	loc.UniformMatrix4f(false, m.Floats())
}
