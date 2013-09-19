package main

import (
	"errors"
	"fmt"
	"github.com/fzipp/geom"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"os"
	"runtime"
)

const (
	resW = 800
	resH = 600
)

func InitGraphics() (*glfw.Window, error) {
	glfw.SetErrorCallback(onError)

	if !glfw.Init() {
		return nil, errors.New("Can't init GLFW.")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	win, err := glfw.CreateWindow(resW, resH, "OpenGL Demo", nil, nil)
	if err != nil {
		return nil, err
	}

	win.MakeContextCurrent()

	// Initialize GLEW
	if gl.Init() != 0 {
		return nil, errors.New("Can't init GLEW.")
	}

	return win, nil
}

func onError(err glfw.ErrorCode, desc string) {
	report(fmt.Errorf("%v: %s", err, desc))
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

const (
	quadVertexCount = 4
	dimensions      = 2
)

type Rect struct {
	geom.Vec2
	geom.Size
}

type quadArray [quadVertexCount * dimensions]float32

func quadCoords(r Rect, coords *quadArray) {
	coords[0] = r.X
	coords[1] = r.Y
	coords[2] = r.X + r.W
	coords[3] = r.Y
	coords[4] = r.X
	coords[5] = r.Y + r.H
	coords[6] = coords[2]
	coords[7] = coords[5]
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	win, err := InitGraphics()
	defer glfw.Terminate()

	if err != nil {
		report(err)
	}

	vs, fs, err := LoadShaders("shaders/vertex.glsl", "shaders/fragment.glsl")
	if err != nil {
		report(err)
	}

	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	var vertices quadArray
	quadCoords(Rect{geom.V2(100, 100), geom.Size{W: 600, H: 400}}, &vertices)

	vao := gl.GenVertexArray()
	vao.Bind()
	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, &vertices, gl.STATIC_DRAW)

	prog, err := LoadProgram(vs, fs)
	if err != nil {
		report(err)
	}
	prog.Use()

	// Init and load PMV matrices
	var pmatrix, mvmatrix geom.Mat4
	pmatrix.Ortho(0, resW, resH, 0, -1.0, 1.0)
	LoadMatrix(&pmatrix, prog, "pmatrix")
	mvmatrix.Id()
	LoadMatrix(&mvmatrix, prog, "mvmatrix")

	attr := prog.GetAttribLocation("vertex")

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		attr.EnableArray()
		vbo.Bind(gl.ARRAY_BUFFER)
		attr.AttribPointer(dimensions, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, quadVertexCount)
		attr.DisableArray()

		win.SwapBuffers()
		glfw.PollEvents()
	}
}
