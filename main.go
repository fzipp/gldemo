package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/fzipp/geom"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	resW = 800
	resH = 600
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func InitGraphics() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(resW, resH, "OpenGL Demo", nil, nil)
	if err != nil {
		return nil, err
	}

	win.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		return nil, err
	}

	return win, nil
}

func ck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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

	ck(err)

	vs, fs, err := LoadShaders("shaders/vertex.glsl", "shaders/fragment.glsl")
	ck(err)

	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	var vertices quadArray
	quadCoords(Rect{geom.V2(100, 100), geom.Size{W: 600, H: 400}}, &vertices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	prog, err := LoadProgram(vs, fs)
	ck(err)
	gl.UseProgram(prog)

	// Init and load PMV matrices
	var pmatrix, mvmatrix geom.Mat4
	pmatrix.Ortho(0, resW, resH, 0, -1.0, 1.0)
	LoadMatrix(&pmatrix, prog, "pmatrix")
	mvmatrix.ID()
	LoadMatrix(&mvmatrix, prog, "mvmatrix")

	attr := uint32(gl.GetAttribLocation(prog, gl.Str("vertex\x00")))

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.EnableVertexAttribArray(attr)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.VertexAttribPointer(attr, dimensions, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, quadVertexCount)
		gl.DisableVertexAttribArray(attr)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}
