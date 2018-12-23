package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unsafe"

	"github.com/fzipp/geom"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func LoadShader(filename string, shaderType uint32) (shader uint32, err error) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	code = append(code, 0)
	shader = gl.CreateShader(shaderType)
	csources, free := gl.Strs(string(code))
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		return shader, fmt.Errorf("failed to compile GL shader: %v", getShaderInfoLog(shader))
	}
	return shader, nil
}

func LoadShaders(vertFile, fragFile string) (vs uint32, fs uint32, err error) {
	vs, err = LoadShader(vertFile, gl.VERTEX_SHADER)
	if err != nil {
		return vs, 0, err
	}
	fs, err = LoadShader(fragFile, gl.FRAGMENT_SHADER)
	if err != nil {
		return vs, fs, err
	}
	return vs, fs, nil
}

func LoadProgram(shaders ...uint32) (program uint32, err error) {
	p := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(p, shader)
	}
	gl.LinkProgram(p)
	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		return p, fmt.Errorf("failed to link GL program: %v", getProgramInfoLog(p))
	}
	gl.ValidateProgram(p)
	gl.GetProgramiv(p, gl.VALIDATE_STATUS, &status)
	if status == gl.FALSE {
		return p, fmt.Errorf("failed to validate GL program: %v", getProgramInfoLog(p))
	}

	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return p, nil
}

func getProgramInfoLog(program uint32) string {
	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return log
}

func getShaderInfoLog(program uint32) string {
	var logLength int32
	gl.GetShaderiv(program, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(program, logLength, nil, gl.Str(log))
	return log
}

func LoadMatrix(m *geom.Mat4, program uint32, name string) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(loc, 1, false, (*float32)(unsafe.Pointer(m.Floats())))
}
